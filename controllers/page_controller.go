package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// returns the home page
func HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func ChallengeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "challenge.html", nil)
}

func AnnouncementsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "announcements.html", nil)
}
func SigninHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "signin.html", nil)
}

func AboutHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", nil)
}

func ContactHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "contact.html", nil)
}

func ClubHomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "clubhome.html", nil)
}

func AdminHomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "adminhome.html", nil)
}

func EditMembersHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "editmembers.html", nil)
}

func ChallengeRulesHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "challengerules.html", nil)
}

func SendAnnouncementsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "sendannouncements.html", nil)
}

func EditPublicPagesHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "editpublicpages.html", nil)
}

func EditMatchesHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "editmatches.html", nil)
}

func ErrorHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "error.html", nil)
}
