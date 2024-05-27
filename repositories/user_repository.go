package repositories

import (
	"ClubTennis/models"
	"errors"
	"sort"

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
	var ranks []uint
	r.db.Model(&models.User{}).Order("`rank` desc").Pluck("`rank`", &ranks).Limit(1)
	if len(ranks) > 0 {
		u.Rank = ranks[0] + 1
	} else {
		u.Rank = 1
	}

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
	return r.db.Save(u).Error
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

// returns the user with the given email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.Where(&models.User{Email: email}).Take(&u).Error
	return &u, err
}

type UserMatch struct {
	MatchID uint
	UserID  uint
}

func (r *UserRepository) DeleteByID(id uint) error {
	var userMatches []UserMatch
	r.db.Table("user_matches").Unscoped().Where("user_id = ?", id).Delete(&userMatches)

	return r.db.Unscoped().Delete(&models.User{}, id).Error
}

func (r *UserRepository) DeleteByUnityID(unityID string) error {
	user := models.User{UnityID: unityID}
	r.db.Select("id").Where(user).Take(&user)
	return r.DeleteByID(user.ID)
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

func (r *UserRepository) FixLadder() {
	ladder, err := r.FindAll()
	if err != nil {
		return
	}
	sort.Slice(ladder, func(i, j int) bool {
		return ladder[i].Rank < ladder[j].Rank
	})

	for i := range ladder {
		ladder[i].Rank = uint(i + 1)
	}
	r.SaveUsers(ladder)
}
