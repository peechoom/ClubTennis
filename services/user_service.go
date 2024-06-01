package services

import (
	"ClubTennis/models"
	"ClubTennis/repositories"
	"errors"
	"sort"

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
		if u[0].ID == 0 {
			return s.repo.SubmitUser(u[0])
		}
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

// returns all users between rank a and b (inclusive) for the given ladder
func (s *UserService) FindByRankRange(ladder string, a uint, b uint) (u []User, err error) {
	if a > b {
		u, err = s.repo.FindByRankRange(b, a, ladder)
	} else {
		u, err = s.repo.FindByRankRange(a, b, ladder)
	}
	return
}

// returns the user with the given rank
func (s *UserService) FindByRank(r uint) (*User, error) {
	return s.repo.FindByRank(r)
}

// returns the user with the given ID
func (s *UserService) FindByID(r uint) (*User, error) {
	return s.repo.FindByID(r)
}

// returns the user with the given unity ID
func (s *UserService) FindByUnityID(unityID string) (*User, error) {
	return s.repo.FindByUnityID(unityID)
}

// returns the user with the given email. It is considered an error if no match is found
func (s *UserService) FindByEmail(email string) (*User, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil || u.ID == 0 {
		return nil, err
	}
	return u, nil
}

// should be used mostly for testing. returns all users in the db
func (s *UserService) FindAll() (u []User, err error) {
	return s.repo.FindAll()
}

// deletes the user with the given numeric id
func (s *UserService) DeleteByID(id uint) error {
	err := s.repo.DeleteByID(id)
	if err != nil {
		return err
	}
	s.repo.FixLadder()
	return nil
}

func (s *UserService) DeleteByUnityID(unityID string) error {
	err := s.repo.DeleteByUnityID(unityID)
	if err != nil {
		return err
	}
	s.repo.FixLadder()
	return nil
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
	if winner.Ladder != loser.Ladder {
		return errors.New("cannot resolve, players are not in the same ladder")
	}

	var ladder []User
	ladder, err := s.FindByRankRange(winner.Ladder, winner.Rank, loser.Rank)

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
	sort.Slice(ladder, func(i, j int) bool {
		return ladder[i].Rank < ladder[j].Rank
	})
	size := len(ladder)

	ladder[size-1].Rank = ladder[0].Rank

	for i := 0; i < size-1; i++ {
		ladder[i].Rank++
	}
}
