package services

import (
	"ClubTennis/models"

	"gorm.io/gorm"
)

type ResetService struct {
	db *gorm.DB
}

func NewResetService(db *gorm.DB) *ResetService {
	return &ResetService{db: db}
}

func (s *ResetService) DeleteEverything() {
	s.db.Migrator().DropTable(&models.User{}, &models.Match{}, &models.Image{}, &models.Announcement{}, &models.Snippet{})
	err := s.db.AutoMigrate(&models.User{}, &models.Match{}, &models.Image{}, &models.Announcement{}, &models.Snippet{})
	if err != nil {
		panic(err.Error())
	}
}
