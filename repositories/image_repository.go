package repositories

import (
	"ClubTennis/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

// save an image to the db. does nothing if filename already exists because 99.9999% chance its the same image
func (r *ImageRepository) SaveImage(img *models.Image) error {
	return r.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "file_name"}},
			DoNothing: true,
		},
	).Create(img).Error
}

func (r *ImageRepository) FindByFileName(filename string) (*models.Image, error) {
	var m models.Image
	err := r.db.Where(&models.Image{FileName: filename}).Take(&m).Error
	return &m, err
}

func (r *ImageRepository) DeleteByFileName(filename string) error {
	var img models.Image
	err := r.db.Select("id", "file_name").Where(&models.Image{FileName: filename}).Take(&img).Error
	if err != nil {
		return err
	}
	return r.db.Unscoped().Delete(&img).Error
}
