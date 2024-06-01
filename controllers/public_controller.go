package controllers

import (
	"ClubTennis/models"
	"ClubTennis/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PublicController struct {
	publicService *services.PublicService
}

func NewPublicController(publicService *services.PublicService) *PublicController {
	if publicService == nil {
		return nil
	}
	return &PublicController{publicService: publicService}
}

// --------------------------------------------------------------------------------------
// GET routings

/*
	GET .../slides

gets the slideshow for the homepage
*/
func (p *PublicController) GetSlideshow(c *gin.Context) {
	slides, err := p.publicService.GetSlideshow()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get slides"})
		return
	}
	c.JSON(http.StatusOK, slides)
}

/*
	GET .../welcome

gets the custom welcome snippet for the homepage
*/
func (p *PublicController) GetWelcome(c *gin.Context) {
	s, err := p.publicService.GetCustomHomePage()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load custom welcome page"})
		return
	}
	//maybe this should be json idk
	c.JSON(http.StatusOK, s)
}

// --------------------------------------------------------------------------------------
// PUT routings
/*
	PUT /admin/slides/:slideNum

uploads new slides to the homepage. expects a slidenum in the url and json with a data field containing the image
base64 representation
*/
func (p *PublicController) PutSlides(c *gin.Context) {
	str := c.Param("slideNum")

	slideNum, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slide num not valid"})
		return
	}
	var payload gin.H
	c.BindJSON(&payload)
	err = p.publicService.PutSlide(int(slideNum), (payload["data"].(string)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slide not allowed. is it too big?"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "accepted"})
}

/*
	PUT /admin/welcome

changes the custom greeting snippet on the welcome page to the html provided
*/
func (p *PublicController) PutWelcome(c *gin.Context) {
	var payload gin.H
	c.BindJSON(&payload)

	if payload["data"] == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data field must contain html"})
		return
	}

	snippet := models.NewSnippet("", payload["data"].(string))
	if snippet == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "illegal fields in html"})
		return
	}
	err := p.publicService.SetCustomHomePage(snippet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save page"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "accepted"})
}
