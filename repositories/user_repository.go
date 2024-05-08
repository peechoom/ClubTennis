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
	err := r.db.Preload("Matches").First(&u, ID).Error

	if err != nil || u.ID != ID || u.UnityID == "" {
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

// submits a slice of new users to the DB, must be USER STRUCTS, not pointers
func (r *UserRepository) SubmitUsers(u []models.User) error {
	err := r.db.Create(&u).Error
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

// updates several user's info in the db. must be a slice of users, NOT user pointers
func (r *UserRepository) SaveUsers(u []models.User) (linesUpdated int64, err error) {
	t := r.db.Save(&u)
	err = t.Error
	linesUpdated = t.RowsAffected
	if err != nil {
		return linesUpdated, err
	}
	return linesUpdated, err
}

// finds the user with the given rank
func (r *UserRepository) FindByRank(Rank uint) (*models.User, error) {
	var u models.User
	if err := r.db.Preload("Matches").Where(&models.User{Rank: Rank}).Take(&u).Error; err != nil {
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
	err := r.db.Preload("Matches").Where("`rank` BETWEEN ? AND ?", LoRank, HiRank).Find(&u).Error

	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindAll() (u []models.User, err error) {
	err = r.db.Preload("Matches").Find(&u).Error
	return
}

func (r *UserRepository) FindByUnityID(s string) (u *models.User, err error) {
	u = new(models.User)
	err = r.db.Preload("Matches").Where(&models.User{UnityID: s}).Find(u).Error

	if u.UnityID != s {
		return nil, nil
	}
	return
}
