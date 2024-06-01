package services

import (
	"ClubTennis/models"
	"ClubTennis/repositories"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"gorm.io/gorm"
)

// service for all public-facing resources
type PublicService struct {
	slideRepo   *repositories.SlideRepository
	snippetRepo *repositories.SnippetRepository
}

// max of 10 mb
const MAX_IMG_SIZE int = 10 * 1024 * 1024

func NewPublicService(db *gorm.DB) *PublicService {
	s := PublicService{slideRepo: repositories.NewSlideRepository(db), snippetRepo: repositories.NewSnippetRepository(db)}
	if (&s).ensureSlideCount() != nil {
		return nil
	}
	return &s
}

// ensures there are 5 slides. fills in missing slides with nc state wallpaper
func (s *PublicService) ensureSlideCount() error {
	slides, err := s.GetSlideshow()
	if err != nil {
		return err
	}

	missing := make([]int, repositories.SLIDE_COUNT)
	for _, s := range slides {
		missing[s.SlideNum] = 1
	}
	for i, n := range missing {
		if n == 0 {
			content, err := encodeImageToBase64("./static/wallpaper.jpg")
			if err != nil {
				return err
			}
			s.PutSlide(i, content)
		}
	}

	return nil
}

func (s *PublicService) GetSlideshow() ([]models.Slide, error) {
	return s.slideRepo.GetSlideshow()
}

func (s *PublicService) PutSlide(slideNum int, slide string) error {
	l := len(slide)
	if l <= 0 || l > MAX_IMG_SIZE {
		return errors.New("size out of bounds")
	}
	return s.slideRepo.SaveSlide(slideNum, slide)
}

func (s *PublicService) SetCustomHomePage(snippet *models.Snippet) error {
	return s.snippetRepo.SetCustomHomePage(snippet)
}

func (s *PublicService) GetCustomHomePage() (*models.Snippet, error) {
	snip := s.snippetRepo.GetCustomHomePage()
	if snip == nil {
		return nil, errors.New("homepage not yet set")
	}
	return snip, nil
}

func encodeImageToBase64(filepath string) (string, error) {
	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read the image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Encode the image to JPEG format
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, nil)
	if err != nil {
		return "", err
	}

	// Convert to Base64 encoding
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Prefix with data URI scheme
	dataURI := fmt.Sprintf("data:image/jpg;base64,%s", base64Str)

	return dataURI, nil
}
