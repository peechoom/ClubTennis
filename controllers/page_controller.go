package controllers

import (
	"ClubTennis/services"
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// only needed for pages that do htmx
type PageController struct {
	snippetservice *services.SnippetService
}

func NewPageController(snippetservice *services.SnippetService) *PageController {
	return &PageController{snippetservice: snippetservice}
}

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

func SendAnnouncementsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "sendannouncements.html", nil)
}

func EditPublicPagesHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "editpublicpages.html", nil)
}

func EditMatchesHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "editmatches.html", nil)
}

func EditRulesHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "editrules.html", nil)
}

func HowToHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "howto.html", nil)
}

func WipeServerHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "wipeserver.html", gin.H{"root": os.Getenv("EMAIL_USERNAME")})
}

func (ctrl *PageController) LadderRulesHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "rules_template.html", gin.H{
		"snippet": template.HTML(ctrl.snippetservice.Get(services.LADDER_CATEGORY).Data),
	})
}

func (ctrl *PageController) ChallengeRulesHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "rules_template.html", gin.H{
		"snippet": template.HTML(ctrl.snippetservice.Get(services.CHALLENGE_CATEGORY).Data),
	})
}

func ErrorHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "error.html", nil)
}
