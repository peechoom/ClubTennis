package repositories

import (
	"ClubTennis/models"
	"errors"

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

// submits new match to the db, updates passed match to hold primary key.
// returns an error if the match is marked inactive or has a score
func (r *MatchRepository) SubmitMatch(m *Match) error {
	if m.Score != 0 || !m.IsActive {
		return errors.New("match must be open to submit")
	}
	err := r.db.Create(m).Error
	if err != nil {
		return err
	}
	return nil
}

// updates an already-existing match in the db. Cascade updates player scores
func (r *MatchRepository) SaveMatch(m *Match) error {
	if m.Score != 0 {
		ur := NewUserRepository(r.db)
		if err := ur.SaveUser(m.Challenger()); err != nil {
			return err
		}
		if err := ur.SaveUser(m.Challenged()); err != nil {
			return err
		}
	}
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

func (r *MatchRepository) FindAll() (m []Match, err error) {
	err = r.db.Preload("Players").Find(&m).Error
	return
}

// returns all matches with the given user ID as the challenger
func (r *MatchRepository) FindByChallengerID(challengerID uint) (m []Match, err error) {
	err = r.db.Preload("Players").Where(&Match{ChallengerID: challengerID}).Find(&m).Error
	return
}

// returns all matches with the given user ID as the challenged player
func (r *MatchRepository) FindByChallengedID(challengedID uint) (m []Match, err error) {
	err = r.db.Preload("Players").Where(&Match{ChallengedID: challengedID}).Find(&m).Error
	return
}

func (r *MatchRepository) FindByPlayerIDs(IDs ...uint) (m []Match, err error) {
	err = r.db.Preload("Players").
		Where("`challenger_id` IN ?", IDs).
		Or("`challenged_id` IN ?", IDs).
		Find(&m).Error
	return
}

func (r *MatchRepository) FindByChallengerChallenged(challengerID uint, challengedID uint) (m []Match, err error) {
	err = r.db.Preload("Players").
		Where(&Match{ChallengerID: challengerID, ChallengedID: challengedID}).
		Find(&m).Error
	return
}
