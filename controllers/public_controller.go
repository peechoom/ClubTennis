package controllers

import (
	"ClubTennis/models"
	"ClubTennis/services"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PublicController struct {
	publicService *services.PublicService
	imageService  *services.ImageService
}

func NewPublicController(publicService *services.PublicService, imageservice *services.ImageService) *PublicController {
	if publicService == nil {
		return nil
	}
	return &PublicController{publicService: publicService, imageService: imageservice}
}

// --------------------------------------------------------------------------------------
// GET routings

/*
	GET .../welcome

gets the custom welcome snippet for the homepage.
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

/*
	GET .../images/{filename}

gets an image from the database. static images are handled with a request to /static/{filename}
*/
func (p *PublicController) GetImage(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}
	img := p.imageService.Get(filename)
	if img == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	c.Data(http.StatusOK, "image/"+img.Extension, img.Data)
}

// --------------------------------------------------------------------------------------
// PUT/POST routings

/*
	POST /admin/slides/:slideNum

uploads new slides to the homepage. expects a slidenum in the url and json with a data field containing the image
base64 representation
*/
func (p *PublicController) PostSlides(c *gin.Context) {
	slideStr := c.Param("slideNum")
	num, err := strconv.ParseInt(slideStr, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "slide number not present/correct"})
		return
	}
	if int(num) > services.SLIDE_COUNT || int(num) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "slide number out of range"})
		return
	}

	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "image data not present"})
		return
	}
	if strings.Contains(file.Filename, "/") {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad filename"})
		return
	}
	// file should be less than 8MB
	if file.Size > 8*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "image too big"})
		return
	}
	reg := regexp.MustCompile(`^.*\.(?:png|jpg|jpeg|webp|PNG|JPG|JPEG|WEBP)$`)
	if !reg.MatchString(file.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad filetype"})
		return
	}

	oldFilename := os.TempDir() + "/" + file.Filename
	c.SaveUploadedFile(file, oldFilename)
	defer os.Remove(oldFilename)

	err = services.ConvertToWebp(oldFilename, "static/slide"+slideStr+".webp")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error converting to webp"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file uploaded successfully"})

	go func(from, to string) {
		var in, out *os.File
		in, _ = os.Open(from)
		defer in.Close()
		out, _ = os.Create(to)
		defer out.Close()
		_, err = io.Copy(out, in)
		return
	}("static/slide"+slideStr+".webp", os.Getenv("SERVER_FILES_MOUNTPOINT")+"/slide"+slideStr+".webp")
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
