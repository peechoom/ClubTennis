package services

import (
	"ClubTennis/models"
	"ClubTennis/repositories"
	"errors"

	"gorm.io/gorm"
)

type MatchService struct {
	repo     *repositories.MatchRepository
	userRepo *repositories.UserRepository
}

// how long we should consider "recent matches"
const recentMatchesDays = 7

type Match = models.Match

func NewMatchService(db *gorm.DB) *MatchService {
	return &MatchService{
		repo:     repositories.NewMatchRepository(db),
		userRepo: repositories.NewUserRepository(db)}
}

// saves a match to the database, creating it if it DNE already, and updating it if it exists
func (ms *MatchService) Save(match *Match) error {
	if match.ID == 0 {
		return ms.repo.SubmitMatch(match)
	}

	fetched, _ := ms.repo.FindByID(match.ID)
	if fetched.ID != match.ID || fetched.ChallengedID != match.ChallengedID || fetched.ChallengerID != match.ChallengerID {
		return errors.New("provided match ID does not coorelate to correct match in database")
	}

	err := ms.repo.SaveMatch(match)
	if err != nil {
		return err
	}

	err = ms.adjustLadder(match)
	return err //hopefully nil
}

func (ms *MatchService) adjustLadder(m *Match) error {
	us := &UserService{repo: ms.userRepo}
	winnerID, c := m.Winner()
	winner, err := us.FindByID(winnerID)
	if err != nil {
		return err
	}
	var loserID uint
	if c {
		loserID = m.ChallengedID
	} else {
		loserID = m.ChallengerID
	}

	loser, err := us.FindByID(loserID)
	if err != nil {
		return err
	}

	return us.AdjustLadder(winner, loser)
}

// finds all matches involving the club members with the given IDs. can be one or several
func (ms *MatchService) FindByPlayerID(IDs ...uint) (m []Match, err error) {
	if len(IDs) == 0 {
		return nil, nil
	}
	m, err = ms.repo.FindByPlayerIDs(IDs...)
	return
}

// finds all matches where the given user (id) was a challenger
func (ms *MatchService) FindByChallengerID(ID uint) (m []Match, err error) {
	m, err = ms.repo.FindByChallengerID(ID)
	return
}

// returns all records in the database
func (ms *MatchService) FindAll() (m []Match, err error) {
	m, err = ms.repo.FindAll()
	return
}

// returns the match with the given id
func (ms *MatchService) FindByID(ID uint) (m *Match, err error) {
	m, err = ms.repo.FindByID(ID)
	return
}

// returns all matches that meet the query
func (ms *MatchService) FindByPlayerIDAndActive(active bool, IDs ...uint) (m []Match, err error) {
	if len(IDs) == 0 {
		return nil, nil
	}
	m, err = ms.repo.FindByPlayerIDAndActive(active, IDs...)
	return
}

func (ms *MatchService) FindByActive(active bool) (m []Match, err error) {
	m, err = ms.repo.FindByActive(active)
	return
}

// TODO implement this
// func (ms *MatchService) FindAllRecentMatches(timespan int) (m []Match, err error) {
// 	m, err =
// }
