package controllers

import (
	"ClubTennis/models"
	"ClubTennis/services"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AnnouncementController struct {
	announcementService *services.AnnouncementService
	userService         *services.UserService
	emailService        *services.EmailService
	imageService        *services.ImageService
	serverHost          string //this servers host. for stripping images and replacing with links
}

func NewAnnouncementController(announcementService *services.AnnouncementService, emailService *services.EmailService, userService *services.UserService, imageService *services.ImageService) *AnnouncementController {
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	if host == "" || port == "" {
		log.Fatal("host and/or port not specified in .env file")
	}
	return &AnnouncementController{
		announcementService: announcementService,
		emailService:        emailService,
		userService:         userService,
		imageService:        imageService,
		serverHost:          "http://" + host}
}

//---------------------------------------------------------------------------------------------------------
// GET HANDLERS

/*
	GET .../announcements/:page

returns page n of the announcements. 0-indexed
*/

func (a *AnnouncementController) GetAnnouncementPage(c *gin.Context) {
	pgString := c.Param("page")
	page, err := strconv.ParseUint(pgString, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page is not a number"})
		return
	}

	ann, err := a.announcementService.GetAnnouncementPage(int(page))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cant get that page"})
		return
	}

	c.JSON(http.StatusOK, ann)
}

// -----------------------------------------------------------------------------------------------------------
// POST HANDLERS

/*
	POST .../announcements

uploads a new announcement to the db
*/
func (a *AnnouncementController) SubmitPost(c *gin.Context) {
	var payload gin.H
	c.BindJSON(&payload)
	if payload["body"] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "body missing"})
	}
	if payload["subject"] == nil {
		payload["subject"] = ""
	}

	ann := models.NewAnnouncement((payload["body"].(string)), (payload["subject"].(string)))
	if ann == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "announcement not valid"})
		return
	}

	if ann.Size() > 7 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Images are too large! Please try compressing or reducing the size of your images"})
		return
	}

	//strip image base64's from announcement and replace them with links that trigger GET's from image repo
	images, err := models.StripImages(&ann.Data, a.serverHost+"/images")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not parse html doc"})
		return
	}

	for _, img := range images {
		err := a.imageService.Save(img)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "could not save image"})
			return
		}
	}

	err = a.announcementService.SubmitAnnouncement(ann)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save announcement"})
		return
	}

	if payload["notifyAll"] != nil && payload["notifyAll"].(bool) {
		err = a.emailEveryone(c, ann)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})

}

func (a *AnnouncementController) emailEveryone(c *gin.Context, ann *models.Announcement) error {
	everyone, err := a.userService.FindAll()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could fetch recipients from database"})
		return err
	}
	emails := a.emailService.MakeAnnouncementEmail(ann, everyone)

	if emails == nil {
		err = errors.New("could not create email arr")
		return err
	}
	for i, email := range emails {
		if email == nil {
			return errors.New("could not send email " + strconv.FormatInt(int64(i), 10))
		}
	}
	var errs []error
	found := false
	for _, email := range emails {
		err := a.emailService.Send(email)
		if err != nil {
			found = true
			errs = append(errs, err)
		}
	}
	if found {
		str := ""
		for _, err := range errs {
			str += err.Error() + "\n"
		}
		return errors.New("The following errors were reported: " + str)
	}

	return nil
}

// -------------------------------------------------------------------------------------------
// DELETE handler
/*

	.../announcements/:id

*/
func (a *AnnouncementController) DeleteAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.String(http.StatusBadRequest, "id field not number/present")
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		c.String(http.StatusBadRequest, "id field not number/present")
		return
	}

	ann, err := a.announcementService.GetAnnouncementByID(uint(id))
	if err != nil || ann.ID == 0 {
		c.Error(err)
		c.String(http.StatusNotFound, "announcement not found")
		return
	}
	go a.deleteImages(strings.Clone(ann.Data))
	a.announcementService.DeleteAnnouncement(ann.ID)

	c.JSON(http.StatusOK, gin.H{"message": "announcement deleted successfully"})
}

func (a *AnnouncementController) deleteImages(body string) {
	//thanks chatgpt
	re := regexp.MustCompile(`<img[^>]*\bsrc=["'][^>]*([a-fA-F0-9]{32}\.\w{1,5})["'][^>]*>`)
	matches := re.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		println("deleting " + match[1])
		err := a.imageService.Delete(match[1])
		if err != nil {
			println(err)
		}
	}
}
