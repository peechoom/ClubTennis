package models

import (
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"gorm.io/gorm"
)

type Match struct {
	gorm.Model               //DO NOT touch the ID, the database handles this
	SubmittedAt    time.Time `gorm:"default:0"` //When this match was marked done. only relevant if IsActive is false
	ChallengerID   uint      //ID of player X who challenged player Y to this match
	ChallengedID   uint      //ID of player Y who challenged player X to this match
	ChallengerRank uint      //rank of the challenger
	ChallengedRank uint      //rank of the challenged player
	Players        []*User   `gorm:"many2many:user_matches;"`
	Score          uint8     `gorm:"default:0"`
	IsActive       bool
	LateNotifSent  bool //if this challenge is almost expired, has the user been notified yet?
}

// A player may not challenge the same player within this many days
const SamePlayerCooldownDays = 14

// A player may not challenge the same player within this many hours
const SamePlayerCooldownHours int = 24 * SamePlayerCooldownDays

// winners score
const WinningScore uint = 6

func (m *Match) Challenger() (u *User) {
	for _, usr := range m.Players {
		if usr.ID == m.ChallengerID {
			return usr
		}
	}
	return nil
}
func (m *Match) Challenged() (u *User) {
	for _, usr := range m.Players {
		if usr.ID == m.ChallengedID {
			return usr
		}
	}
	return nil
}

// returns the winner's id and wether or not the challenger won
func (m *Match) Winner() (uint, bool) {
	if m.Score == 0 {
		return 0, false
	}

	a, _ := DecodeScore(m.Score)
	if a == int(WinningScore) {
		return m.ChallengerID, true
	}
	return m.ChallengedID, false
}

// constructor for a new match
func NewMatch(Challenger *User, Challenged *User) *Match {
	if Challenger == Challenged {
		return nil
	}
	m := new(Match)
	m.Players = append(m.Players, Challenger)
	m.Players = append(m.Players, Challenged)

	m.ChallengerID = Challenger.ID
	m.ChallengedID = Challenged.ID

	m.ChallengerRank = Challenger.Rank
	m.ChallengedRank = Challenged.Rank

	m.IsActive = true
	log.Print(m.SubmittedAt)

	return m
}

func (challenger *User) CanChallenge(challenged *User) (bool, error) {
	if challenger.Ladder != challenged.Ladder {
		return false, errors.New("cannot challenge players in other ladders")
	}
	if challenger.UnityID == challenged.UnityID {
		return false, errors.New("players cannot challenge themselves")
	}
	if challenger.HasActiveMatch() {
		return false, errors.New("challenger has an active challenge")
	}
	if challenged.HasActiveMatch() {
		return false, errors.New("challenged player already has an active challenge")
	}
	if !challenger.IsWithinRangeOf(challenged) {
		return false, errors.New("challenger rank is too low")
	}
	if challenger.HasRecentlyChallenged(challenged) {
		return false, fmt.Errorf("challenger has challenged player within the last %d days", SamePlayerCooldownDays)
	}
	if !challenger.IsActive {
		return false, fmt.Errorf("Challenger is not active")
	}
	if !challenged.IsActive {
		return false, fmt.Errorf("Challenged player is not active")
	}

	return true, nil
}

/*
function for a user to challenge another user to a match.
returns a pointer to the created match on success, or nil + error if error
*/
func (challenger *User) Challenge(challenged *User) (*Match, error) {
	c, err := challenger.CanChallenge(challenged)
	if !c {
		return nil, err
	}

	var match *Match = NewMatch(challenger, challenged)

	challenger.Matches = append(challenger.Matches, match)
	challenged.Matches = append(challenged.Matches, match)

	return match, nil
}

// returns true if the user has an active match that has not been played yet
func (u *User) HasActiveMatch() bool {
	for _, m := range u.Matches {
		if m.IsActive {
			return true
		}
	}
	return false
}

// returns true if the challenger is high enough in the ladder to challenge the user.
// A players challenging range decreases as their seed gets closer to 1.
func (challenger *User) IsWithinRangeOf(challenged *User) bool {
	var diff int = int(challenger.Rank) - int(challenged.Rank)

	switch true {
	case challenger.Rank <= 10:
		return diff <= 2
	case challenger.Rank < 20:
		return diff <= 3
	case challenger.Rank <= 35:
		return diff <= 5
	}
	return diff <= 8
}

// returns true if the challenger has challenged the challenged player  within the last 14 days.
func (challenger *User) HasRecentlyChallenged(challenged *User) bool {
	for _, match := range challenger.Matches {
		if match.ChallengedID == challenged.ID && (match.IsActive || time.Since(match.CreatedAt).Hours() < float64(SamePlayerCooldownHours)) {
			return true
		}
	}
	return false
}

func (match *Match) SubmitScore(challengerScore int, challengedScore int) (err error) {
	if challengerScore > 6 || challengedScore > 6 || challengerScore < 0 || challengedScore < 0 {
		return errors.New("illegal score")
	}

	if (challengedScore != 6 && challengerScore != 6) || (challengerScore == 6 && challengedScore == 6) {
		return errors.New("illegal score")
	}

	if challengerScore == 6 {
		match.Challenger().Wins++
		match.Challenged().Losses++
	} else {
		match.Challenged().Wins++
		match.Challenger().Losses++
	}

	match.SubmittedAt = time.Now()
	match.Score = EncodeScore(challengerScore, challengedScore)
	match.IsActive = false

	return nil
}

func (match *Match) Cancel() {
	match.IsActive = false
	match.SubmittedAt = time.Now()
	match.Score = 0
	for _, p := range match.Players {
		removeFromMatches(p, match)
	}
}

// encodes a valid score into an unsigned 8 bit number
func EncodeScore(challengerScore int, challengedScore int) uint8 {
	return (uint8(challengerScore) << 4) | uint8(challengedScore)
}

// decodes the score from an unsigned 8 bit number
func DecodeScore(score uint8) (challengerScore int, challengedScore int) {
	return int(score >> 4), int(score & 0xF)
}

func removeFromMatches(user *User, match *Match) {
	if user == nil || match == nil {
		return
	}
	user.Matches = slices.DeleteFunc(user.Matches, func(m *Match) bool { return m.ID == match.ID })
}
