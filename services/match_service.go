package services

import (
	"ClubTennis/models"
	"ClubTennis/repositories"
	"errors"
	"time"

	"gorm.io/gorm"
)

type MatchService struct {
	repo     *repositories.MatchRepository
	userRepo *repositories.UserRepository
}

// how long we should consider "recent matches"
// const recentMatchesDays = 7

const matchExpiresDays = 9

const deletionThreshold = models.SamePlayerCooldownDays * 2

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

// returns all matches that are/arent active
func (ms *MatchService) FindByActive(active bool) (m []Match, err error) {
	m, err = ms.repo.FindByActive(active)
	return
}

// returns all submitted that are at most {timespan} old. if so many matches
// cannot be filled, open matches will fill in the gaps
func (ms *MatchService) FindAllRecentMatches(timespan time.Duration) ([]Match, error) {
	const recent_size int = 15
	m, err := ms.repo.FindByMaxAge(timespan)
	var a []Match
	if err != nil || len(m) < recent_size {
		a, err = ms.repo.FindByActive(true)
		if err != nil {
			return nil, err
		}
	}
	for i := 0; len(m) < recent_size && i < len(a); i++ {
		m = append(m, a[i])
	}
	return m, nil
}

// returns all matches that are about to expire
func (ms *MatchService) FindByNearlyExpired() ([]Match, error) {
	minTime := time.Hour * 24 * (matchExpiresDays - 1)
	maxTime := time.Hour * 24 * (matchExpiresDays)

	return ms.repo.FindByAgeRange(minTime, maxTime, true)
}

func (ms *MatchService) CancelMatch(ID uint) (err error) {
	m, err := ms.FindByID(ID)
	if err != nil || m == nil {
		return err
	}
	m.Cancel()

	var players []models.User
	for _, p := range m.Players {
		players = append(players, *p)
	}

	_, err = ms.userRepo.SaveUsers(players)
	if err != nil {
		return err
	}
	return ms.repo.Delete(*m)
}

// returns all matches that are expired (marked active but too much time has passed)
func (ms *MatchService) FindByExpired() ([]Match, error) {
	minTime := time.Hour * 24 * (matchExpiresDays)
	maxTime := time.Hour * 24 * (matchExpiresDays + 1)

	return ms.repo.FindByAgeRange(minTime, maxTime, true)
}

func (ms *MatchService) DeleteOldMatches() error {
	minTime := time.Hour * 24 * deletionThreshold
	maxTime := time.Hour * 24 * (365)

	matches, err := ms.repo.FindByAgeRange(minTime, maxTime, false)
	if err != nil {
		return err
	}
	return ms.repo.Delete(matches...)
}
