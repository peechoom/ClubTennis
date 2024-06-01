package controllers

import (
	"ClubTennis/models"
	"ClubTennis/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnnouncementController struct {
	announcementService *services.AnnouncementService
}

func NewAnnouncementController(announcementService *services.AnnouncementService) *AnnouncementController {
	return &AnnouncementController{announcementService: announcementService}
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

	ann := models.NewAnnouncement((payload["body"].(string)))
	if ann == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "announcement is not valid HTML"})
		return
	}

	err := a.announcementService.SubmitAnnouncement(ann)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save announcement"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"body": "success"})
	if payload["notifyAll"].(bool) {
		// TODO notify everyone via email
		return
	}
}
