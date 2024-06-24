package repositories

import (
	"ClubTennis/models"
	"errors"

	"gorm.io/gorm"
)

type AnnouncementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) *AnnouncementRepository {
	if db == nil {
		return nil
	}
	return &AnnouncementRepository{db: db}
}

func (r *AnnouncementRepository) SubmitAnnouncement(ann *models.Announcement) error {
	if len(ann.Data) == 0 {
		return errors.New("cannot sumbit a blank announcement")
	}
	return r.db.Create(ann).Error
}

// gets the n'th page of announcements, where there are 'perPage' announcements per page
func (r *AnnouncementRepository) GetAnnouncementPage(page int, perPage int) (ann []models.Announcement, err error) {
	err = r.db.Order("created_at DESC").Limit(perPage).Offset(page * perPage).Find(&ann).Error
	return
}

func (r *AnnouncementRepository) GetAnnouncementByID(ID uint) (*models.Announcement, error) {
	var ann models.Announcement
	ann.ID = ID
	err := r.db.First(&ann).Error
	return &ann, err
}

func (r *AnnouncementRepository) EditAnnouncement(ann *models.Announcement) error {
	if ann.ID == 0 {
		return errors.New("must include announcement id in announcement in order to edit")
	}

	return r.db.Save(ann).Error
}

func (r *AnnouncementRepository) DeleteAnnouncement(ann *models.Announcement) error {
	if ann.ID == 0 {
		return errors.New("must include announcement id in announcement in order to edit")
	}

	return r.db.Unscoped().Where("id = ?", ann.ID).Delete(ann).Error
}
