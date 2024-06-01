package repositories

import (
	"ClubTennis/models"
	"errors"

	"gorm.io/gorm"
)

type SlideRepository struct {
	db *gorm.DB
}

// the number of slides in the slideshow
const SLIDE_COUNT int = 5

func NewSlideRepository(db *gorm.DB) *SlideRepository {
	return &SlideRepository{db: db}
}

func (r *SlideRepository) GetSlideshow() (m []models.Slide, err error) {
	err = r.db.Find(&m).Error
	return
}

// saves slide number x, where x is between 0 and SLIDE_COUNT and slide is base64 encoded
func (r *SlideRepository) SaveSlide(slideNum int, slide string) error {
	if slideNum < 0 || slideNum > SLIDE_COUNT {
		return errors.New("slide out of bounds")
	}
	var s models.Slide
	r.db.Unscoped().Where("`slide_num` = ?", slideNum).Delete(&s)
	s.ID = 0
	s.Data = slide
	s.SlideNum = slideNum
	return r.db.Save(&s).Error
}
