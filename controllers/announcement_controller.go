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
	images, err := ann.StripImages(a.serverHost + "/images")
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
	c.JSON(http.StatusOK, gin.H{"body": "success"})

	if payload["notifyAll"] != nil && payload["notifyAll"].(bool) {
		a.emailEveryone(c, ann)
	}
}

func (a *AnnouncementController) emailEveryone(c *gin.Context, ann *models.Announcement) error {
	everyone, err := a.userService.FindAll()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could fetch recipients from database"})
		return err
	}
	e := a.emailService.MakeAnnouncementEmail(ann, everyone)
	if e == nil {
		err = errors.New("could not create email")
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create email"})
		return err
	}
	err = a.emailService.Send(e)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "email generated but could not be sent"})
		return err
	}
	return nil
}

func (a *AnnouncementController) DeleteAnnouncement(c *gin.Context) {
	var payload gin.H
	c.BindJSON(payload)
	if payload["id"] == nil {
		c.String(http.StatusBadRequest, "id field not present")
		return
	}
	ann, err := a.announcementService.GetAnnouncementByID(payload["id"].(uint))
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
	re := regexp.MustCompile(`<img[^>]*\bsrc=["']([a-fA-F0-9]{32}\.\w{1,5})["'][^>]*>`)

	matches := re.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		a.imageService.Delete(match[1])
	}
}
