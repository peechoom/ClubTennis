package repositories

import (
	"ClubTennis/models"

	"gorm.io/gorm"
)

type MatchRepository struct {
	db *gorm.DB // pointer to the gorm database
}
type Match = models.Match

// initialize the match repository
func NewMatchRepository(db *gorm.DB) *MatchRepository {
	if db == nil {
		return nil
	}
	return &MatchRepository{db: db}
}

// finds the match with the given ID
func (r *MatchRepository) FindByID(ID uint) (*models.Match, error) {
	var m Match
	err := r.db.Preload("Players").First(&m, ID).Error

	if err != nil {
		return nil, err
	}
	return &m, nil
}

// submits new match to the db, updates passed match to hold primary key
func (r *MatchRepository) SubmitMatch(m *Match) error {
	err := r.db.Create(m).Error
	if err != nil {
		return err
	}
	return nil
}

// updates an already-existing match in the db
func (r *MatchRepository) SaveMatch(m *Match) error {
	err := r.db.Save(m).Error
	if err != nil {
		return err
	}
	return nil
}

// finds all matches marked as active
func (r *MatchRepository) FindByActiveTrue() ([]Match, error) {
	var matches []Match
	err := r.db.Preload("Players").Where(&Match{IsActive: true}).Find(&matches).Error

	if err != nil {
		return nil, err
	}

	return matches, nil
}
