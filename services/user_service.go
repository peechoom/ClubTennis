package services

import (
	"ClubTennis/models"
	"ClubTennis/repositories"
	"errors"

	"gorm.io/gorm"
)

type UserService struct {
	repo *repositories.UserRepository
}
type User = models.User

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{repo: repositories.NewUserRepository(db)}
}

// saves user(s) to the database
func (s *UserService) Save(u ...*User) error {
	if len(u) == 0 {
		return nil
	}
	if len(u) == 1 {
		return s.repo.SaveUser(u[0])
	}
	return s.saveMany(u)
}
func (s *UserService) saveMany(users []*User) error {
	var arr []models.User
	for _, u := range users {
		arr = append(arr, *u)
	}
	_, err := s.repo.SaveUsers(arr)
	return err
}

// returns all users between rank a and b (inclusive)
func (s *UserService) FindByRankRange(a uint, b uint) (u []User, err error) {
	if a > b {
		c := a
		a = b
		b = c
	}
	u, err = s.repo.FindByRankRange(a, b)
	return
}

// returns the user with the given rank
func (s *UserService) FindByRank(r uint) (u *User, err error) {
	u, err = s.repo.FindByRank(r)
	return
}

// returns the user with the given ID
func (s *UserService) FindByID(r uint) (u *User, err error) {
	u, err = s.repo.FindByID(r)
	return
}

// returns the user with the given unity ID
func (s *UserService) FindByUnityID(unityID string) (u *User, err error) {
	u, err = s.repo.FindByUnityID(unityID)
	return
}

// should be used mostly for testing. returns all users in the db
func (s *UserService) FindAll() (u []User, err error) {
	return s.repo.FindAll()
}

// algorithm for adjusting the ladder.
// If the winner is a higher rank (lower number), nothing happens.
// If the winner is a lower rank:
//
//	the winner takes the loser's rank
//	the loser's rank is decreased by 1
//	every player ranked between the winner's initial rank and the loser's initial rank is decreased by 1
//	Analogy: taking a block off of the bottom of a jenga tower and putting it on the top
func (s *UserService) AdjustLadder(winner *User, loser *User) error {
	//ladder is already correct
	if winner.Rank < loser.Rank {
		return nil
	}

	var ladder []User
	ladder, err := s.repo.FindByRankRange(loser.Rank, winner.Rank)

	if err != nil {
		return errors.New("rankings not adjusted, " + err.Error())
	}

	ladderAlgo(ladder)

	lines, err := s.repo.SaveUsers(ladder)
	if err != nil {
		return err
	}
	if lines == 0 {
		return errors.New("nothing was changed despite winner having lower rank")
	}
	return nil
}

func ladderAlgo(ladder []User) {
	len := len(ladder)

	ladder[len-1].Rank = ladder[0].Rank

	for i := 0; i < len-1; i++ {
		ladder[i].Rank++
	}
}
