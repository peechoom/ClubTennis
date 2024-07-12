package controllers

import (
	"ClubTennis/models"
	"ClubTennis/services"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type RulesController struct {
	imageService   *services.ImageService
	snippetService *services.SnippetService
	serverhost     string // for stripping images
}

func NewRulesController(imageService *services.ImageService, snippetService *services.SnippetService, serverhost string) *RulesController {
	return &RulesController{imageService: imageService, snippetService: snippetService, serverhost: serverhost}
}

//----------------------------------------------------------------------------------------------------------
// GET handlers

/*
	GET .../ladderrulessnippet

	returns the snippet for the ladder rules page
*/

func (r *RulesController) GetLadderRules(c *gin.Context) {
	r.getGeneric(c, services.LADDER_CATEGORY)
}

/*
	GET .../challengerulessnippet

returns the snippet for the challenge rules page
*/
func (r *RulesController) GetChallengeRules(c *gin.Context) {
	r.getGeneric(c, services.CHALLENGE_CATEGORY)
}

func (r *RulesController) getGeneric(c *gin.Context, category string) {
	s := r.snippetService.Get(category)
	if s == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load custom welcome page"})
		return
	}
	//maybe this should be json idk
	c.JSON(http.StatusOK, s)
}

//---------------------------------------------------------------------------------------------------------
// PUT handlers

/*
	PUT .../ladderrulessnippet

uploads a new ladder rules snippet
*/
func (r *RulesController) PutLadderRules(c *gin.Context) {
	r.putGeneric(c, services.LADDER_CATEGORY)
}

/*
	PUT .../challengerulessnippet

uploads a new challenge rules snippet
*/
func (r *RulesController) PutChallengeRules(c *gin.Context) {
	r.putGeneric(c, services.CHALLENGE_CATEGORY)
}

func (r *RulesController) putGeneric(c *gin.Context, category string) {
	var payload gin.H
	c.BindJSON(&payload)

	if payload["data"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data field must contain html"})
		return
	}

	old := r.snippetService.Get(category)
	if old != nil { //sohuld always be the case, but not necessarily an error
		go r.deleteImages(old.Data)
	}

	snippet := models.NewSnippet("", payload["data"].(string))
	if snippet == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "illegal fields in html"})
		return
	}
	images, err := models.StripImages(&snippet.Data, r.serverhost+"/images")
	if err != nil {
		c.Error(err)
	}

	for _, img := range images {
		if r.imageService.Save(img) != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "error saving image!"})
		}
	}

	err = r.snippetService.Set(category, snippet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save page"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "accepted"})
}

//---------------------------------------------------------------------------------------------------
// UTILITY functions

// parses and deletes all images from an html body
func (r *RulesController) deleteImages(body string) {
	//thanks chatgpt
	re := regexp.MustCompile(`<img[^>]*\bsrc=["'][^>]*([a-fA-F0-9]{32}\.\w{1,5})["'][^>]*>`)
	matches := re.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		println("deleting " + match[1])
		err := r.imageService.Delete(match[1])
		if err != nil {
			println(err)
		}
	}
}
