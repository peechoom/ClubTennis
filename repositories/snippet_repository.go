package repositories

import (
	"ClubTennis/models"

	"gorm.io/gorm"
)

type SnippetRepository struct {
	db *gorm.DB
}

const homepage_category string = "homepage"

func NewSnippetRepository(db *gorm.DB) *SnippetRepository {
	return &SnippetRepository{db: db}
}

func (r *SnippetRepository) GetCustomHomePage() *models.Snippet {
	var s models.Snippet
	if r.db.Where(&models.Snippet{Category: homepage_category}).First(&s).Error != nil {
		return nil
	}
	return &s
}

func (r *SnippetRepository) SetCustomHomePage(snippet *models.Snippet) error {
	var s = r.GetCustomHomePage()
	if s == nil {
		snippet.Category = homepage_category
		return r.db.Save(snippet).Error
	}
	s.Data = snippet.Data
	return r.db.Save(s).Error
}
