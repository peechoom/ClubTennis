package services

import (
	"ClubTennis/models"
	"ClubTennis/repositories"
	"log"

	"gorm.io/gorm"
)

type ImageService struct {
	repo *repositories.ImageRepository
}

func NewImageService(db *gorm.DB) *ImageService {
	return &ImageService{
		repo: repositories.NewImageRepository(db),
	}
}

// get the image from the databse
func (s *ImageService) Get(filename string) *models.Image {
	img, err := s.repo.FindByFileName(filename)
	if err != nil || img.ID == 0 {
		log.Print(err.Error())
		return nil
	}
	return img
}

// save the image to the databse
func (s *ImageService) Save(img *models.Image) error {
	return s.repo.SaveImage(img)
}

// delete the image from the database
func (s *ImageService) Delete(filename string) error {
	return s.repo.DeleteByFileName(filename)
}
