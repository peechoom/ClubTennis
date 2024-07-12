package repositories

import (
	"ClubTennis/models"

	"gorm.io/gorm"
)

type SnippetRepository struct {
	db *gorm.DB
}

func NewSnippetRepository(db *gorm.DB) *SnippetRepository {
	return &SnippetRepository{db: db}
}

func (r *SnippetRepository) FindByCategory(category string) *models.Snippet {
	if category == "" {
		return nil
	}

	var s models.Snippet
	if r.db.Where(&models.Snippet{Category: category}).First(&s).Error != nil {
		return nil
	}
	return &s
}

func (r *SnippetRepository) FindAll() []models.Snippet {
	var s []models.Snippet
	if r.db.Find(&s).Error != nil {
		return nil
	}
	return s
}

func (r *SnippetRepository) Save(category string, snippet *models.Snippet) error {
	var s models.Snippet
	err := r.db.Where(&models.Snippet{Category: category}).First(&s).Error
	if err != nil {
		snippet.Category = category
		return r.db.Save(snippet).Error
	}
	s.Data = snippet.Data
	return r.db.Save(s).Error
}
