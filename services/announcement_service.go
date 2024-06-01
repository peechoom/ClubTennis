package services

import (
	"ClubTennis/models"
	"ClubTennis/repositories"

	"gorm.io/gorm"
)

type AnnouncementService struct {
	repo *repositories.AnnouncementRepository
}

const ANNOUNCEMENTS_PER_PAGE int = 5

func NewAnnouncementService(db *gorm.DB) *AnnouncementService {
	if db == nil {
		return nil
	}
	return &AnnouncementService{repo: repositories.NewAnnouncementRepository(db)}
}

// returns the n'th page of announcements sorted by newest first. page 0 is the first page
func (s *AnnouncementService) GetAnnouncementPage(page int) ([]models.Announcement, error) {
	return s.repo.GetAnnouncementPage(page, ANNOUNCEMENTS_PER_PAGE)
}

func (s *AnnouncementService) SubmitAnnouncement(ann *models.Announcement) error {
	return s.repo.SubmitAnnouncement(ann)
}

func (s *AnnouncementService) EditAnnouncement(ann *models.Announcement) error {
	return s.repo.EditAnnouncement(ann)
}
