package services

import (
	"ClubTennis/models"
	"ClubTennis/repositories"
	"errors"
	"os"
	"sort"
	"strconv"

	"gorm.io/gorm"
)

// the defaut cutoff for red/white team. the lowest possible rank to still be in red
const DEFAULT_CUTOFF_MENS int = 40
const DEFAULT_CUTOFF_WOMENS int = 20
const CUTOFF_FILENAME_MENS string = "mens_cutoff.txt"
const CUTOFF_FILENAME_WOMENS string = "womens_cutoff.txt"

type UserService struct {
	repo         *repositories.UserRepository
	mensCutoff   int
	womensCutoff int
}
type User = models.User

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{repo: repositories.NewUserRepository(db), mensCutoff: 0, womensCutoff: 0}
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
		err := s.repo.SaveUser(u[0])
		s.repo.FixLadder()
		return err
	}
	return s.saveMany(u)
}
func (s *UserService) saveMany(users []*User) error {
	var arr []models.User
	for _, u := range users {
		arr = append(arr, *u)
	}
	_, err := s.repo.SaveUsers(arr)
	s.repo.FixLadder()
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
func (s *UserService) FindBySigninEmail(email string) (*User, error) {
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

func (s *UserService) FindOfficers() ([]models.User, error) {
	return s.repo.FindAdmins()
}

// get the cutoff for red/white team. This is the lowest ranked that someone can have and still be in red
func (s *UserService) GetLadderCutoff(ladder string) int {
	readCutoffFile := func(filename string) int {
		val, err := os.ReadFile(filename)
		if err != nil {
			print(err.Error())
			return -1
		}
		ret, err := strconv.ParseInt(string(val), 10, 0)
		if err != nil {
			print(err.Error())
			return -1
		}
		return int(ret)
	}

	dir := os.Getenv("SERVER_FILES_MOUNTPOINT")
	if dir == "" {
		print("dir dne")
		return -1
	}

	if ladder == models.MENS_LADDER {
		if s.mensCutoff != 0 {
			return s.mensCutoff
		}

		filename := dir + "/" + CUTOFF_FILENAME_MENS
		if _, e := os.Stat(filename); errors.Is(e, os.ErrNotExist) {
			if e = makeCutoffFile(filename, DEFAULT_CUTOFF_MENS); e != nil {
				print(e.Error())
				return -1
			}
		}
		s.mensCutoff = readCutoffFile(filename)
		return s.mensCutoff

	} else if ladder == models.WOMENS_LADDER {
		if s.womensCutoff != 0 {
			return s.womensCutoff
		}

		filename := dir + "/" + CUTOFF_FILENAME_WOMENS
		if _, e := os.Stat(filename); errors.Is(e, os.ErrNotExist) {
			if e = makeCutoffFile(filename, DEFAULT_CUTOFF_WOMENS); e != nil {
				print(e.Error())
				return -1
			}
		}
		s.womensCutoff = readCutoffFile(filename)
		return s.womensCutoff
	}
	return -1
}

// set the cutoff for red/white team. This is the lowest ranked that someone can have and still be in red
func (s *UserService) SetLadderCutoff(ladder string, rank int) error {
	dir := os.Getenv("SERVER_FILES_MOUNTPOINT")
	if dir == "" {
		return errors.New("dir string empty")
	}
	if rank < 0 {
		return errors.New("cutoff cannot be less than 0")
	}

	if ladder == models.MENS_LADDER {
		filename := dir + "/" + CUTOFF_FILENAME_MENS
		if _, e := os.Stat(filename); e == os.ErrNotExist {
			e = makeCutoffFile(filename, rank)
			if e != nil {
				return e
			}
			s.mensCutoff = rank
			return nil
		}
		if e := os.Remove(filename); e != nil {
			return e
		}
		if e := makeCutoffFile(filename, rank); e != nil {
			return e
		}
		s.mensCutoff = rank
		return nil

	} else if ladder == models.WOMENS_LADDER {
		filename := dir + "/" + CUTOFF_FILENAME_WOMENS
		if _, e := os.Stat(filename); e == os.ErrNotExist {
			e = makeCutoffFile(filename, rank)
			if e != nil {
				return e
			}
			s.womensCutoff = rank
			return nil
		}
		if e := os.Remove(filename); e != nil {
			return e
		}
		if e := makeCutoffFile(filename, rank); e != nil {
			return e
		}
		s.womensCutoff = rank
		return nil
	}
	return errors.New("invalid ladder name")
}

func makeCutoffFile(filename string, cutoff int) error {
	val := []byte(strconv.FormatInt(int64(cutoff), 10))
	return os.WriteFile(filename, val, 0777)
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
