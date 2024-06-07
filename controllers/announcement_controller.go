package controllers

import (
	"ClubTennis/models"
	"ClubTennis/services"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnnouncementController struct {
	announcementService *services.AnnouncementService
	userService         *services.UserService
	emailService        *services.EmailService
}

func NewAnnouncementController(announcementService *services.AnnouncementService, emailService *services.EmailService, userService *services.UserService) *AnnouncementController {
	return &AnnouncementController{announcementService: announcementService, emailService: emailService, userService: userService}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "body missing"})
	}
	if payload["subject"] == nil {
		payload["subject"] = ""
	}

	ann := models.NewAnnouncement((payload["body"].(string)), (payload["subject"].(string)))
	if ann == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "announcement not valid"})
		return
	}

	err := a.announcementService.SubmitAnnouncement(ann)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save announcement"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could fetch recipients from database"})
		return err
	}
	e := a.emailService.MakeAnnouncementEmail(ann, everyone)
	if e == nil {
		err = errors.New("could not create email")
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create email"})
		return err
	}
	err = a.emailService.Send(e)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email generated but could not be sent"})
		return err
	}
	return nil
}
