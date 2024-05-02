package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Match struct {
	gorm.Model
	ChallengerID uint
	ChallengedID uint
	Challenger   *User
	Challenged   *User
	Score        uint8 `gorm:"default:0"`
	IsActive     bool
}

// A player may not challenge the same player within this many days
const SamePlayerCooldownDays int = 14

// A player may not challenge the same player within this many hours
const SamePlayerCooldownHours int = 24 * SamePlayerCooldownDays

// constructor for a new match
func NewMatch(Challenger *User, Challenged *User) *Match {
	if Challenger == Challenged {
		return nil
	}
	m := new(Match)
	m.Challenger = Challenger
	m.Challenged = Challenged

	m.ChallengerID = Challenger.ID
	m.ChallengedID = Challenged.ID

	m.IsActive = true

	return m
}

/*
function for a user to challenge another user to a match.
returns a pointer to the created match on success, or nil + error if error
*/
func (challenger *User) Challenge(challenged *User) (*Match, error) {
	if challenger.UnityID == challenged.UnityID {
		return nil, errors.New("players cannot challenge themselves")
	}
	if challenger.HasActiveMatch() {
		return nil, errors.New("challenger has an active challenge")
	}
	if challenged.HasActiveMatch() {
		return nil, errors.New("challenged player already has an active challenge")
	}
	if !challenger.IsWithinRangeOf(challenged) {
		return nil, errors.New("challenger rank is too low")
	}
	if challenger.HasRecentlyChallenged(challenged) {
		return nil, fmt.Errorf("challenger has challenged player within the last %d days", SamePlayerCooldownDays)
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
		if match.Challenged.UnityID == challenged.UnityID && (match.IsActive || time.Since(match.CreatedAt).Hours() < float64(SamePlayerCooldownHours)) {
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
		match.Challenger.Wins++
		match.Challenged.Losses++
	} else {
		match.Challenged.Wins++
		match.Challenger.Losses++
	}

	match.Score = EncodeScore(challengerScore, challengedScore)
	match.IsActive = false

	return nil
}

// encodes a valid score into an unsigned 8 bit number
func EncodeScore(challengerScore int, challengedScore int) uint8 {
	return (uint8(challengerScore) << 4) | uint8(challengedScore)
}

// decodes the score from an unsigned 8 bit number
func DecodeScore(score uint8) (challengerScore int, challengedScore int) {
	return int(score >> 4), int(score & 0xF)
}
