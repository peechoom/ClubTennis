package repositories

import (
	"ClubTennis/models"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// finds the User with the given ID
func (r *UserRepository) FindByID(ID uint) (*models.User, error) {
	var u models.User
	err := r.db.First(&u, ID).Error

	if err != nil {
		return nil, err
	}
	return &u, nil
}

// submits a new User to the DB
func (r *UserRepository) SubmitUser(u *models.User) error {
	err := r.db.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}

// updates a user's info in the DB
func (r *UserRepository) SaveUser(u *models.User) error {
	if err := r.db.Save(u).Error; err != nil {
		return err
	}
	return nil
}

// finds the user with the given rank
func (r *UserRepository) FindByRank(Rank uint) (*models.User, error) {
	var u models.User
	if err := r.db.Where(&models.User{Rank: Rank}).Take(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// returns all users with ranks in a given range, inclusive -> [lo, hi].
// LoRank is the lower NUMBER jackass. we all know 1 is a "higher ranking" but not here
// returns slice of USERS, --> NOT USER POINTERS <--
func (r *UserRepository) FindByRankRange(LoRank uint, HiRank uint) ([]models.User, error) {
	if LoRank > HiRank {
		return nil, errors.New("LoRank must be lower than HiRank")
	}

	var u []models.User
	err := r.db.Where("`rank` BETWEEN ? AND ?", LoRank, HiRank).Find(&u).Error

	if err != nil {
		return nil, err
	}
	return u, nil
}
