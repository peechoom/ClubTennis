package models_test

import (
	"ClubTennis/models"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type User = models.User
type Match = models.Match

func getTestUsers() (*User, *User) {
	var userA *User = models.NewUser("ajhende3", "ncsu", "Alec", "Henderson")
	var userB *User = models.NewUser("alhunt6", "ncsu", "Alison", "Hunt")

	userA.ID = 1
	userB.ID = 2
	return userA, userB
}

func getTestUsersWithRanks(rankA int, rankB int) (*User, *User) {
	userA, userB := getTestUsers()
	userA.Rank = uint(rankA)
	userB.Rank = uint(rankB)
	return userA, userB
}

// fails if the match does not contain the challenger and challenged user
func requireMatchContains(t *testing.T, match *Match, challenger *User, challenged *User) {
	require.Equal(t, challenger.ID, match.ChallengerID)
	require.Equal(t, challenged.ID, match.ChallengedID)
	require.Equal(t, challenger.UnityID, match.Challenger.UnityID)
	require.Equal(t, challenged.UnityID, match.Challenged.UnityID)
}

func TestNew(t *testing.T) {
	userA, userB := getTestUsers()

	var match *models.Match = models.NewMatch(userA, userB)

	requireMatchContains(t, match, userA, userB)
}

func TestChallengeValid(t *testing.T) {
	var userA *User
	var userB *User
	var match *Match
	var err error

	userA, userB = getTestUsersWithRanks(2, 1)
	match, err = userA.Challenge(userB)
	require.NoError(t, err)
	require.NotNil(t, match)
	requireMatchContains(t, match, userA, userB)
	require.Contains(t, userA.Matches, match)
	require.Contains(t, userB.Matches, match)

	userA, userB = getTestUsersWithRanks(3, 1)
	match, err = userA.Challenge(userB)
	require.NoError(t, err)
	require.NotNil(t, match)
	requireMatchContains(t, match, userA, userB)

	userA, userB = getTestUsersWithRanks(10, 8)
	match, err = userA.Challenge(userB)
	require.NoError(t, err)
	require.NotNil(t, match)
	requireMatchContains(t, match, userA, userB)

	userA, userB = getTestUsersWithRanks(11, 8)
	match, err = userA.Challenge(userB)
	require.NoError(t, err)
	require.NotNil(t, match)
	requireMatchContains(t, match, userA, userB)

	userA, userB = getTestUsersWithRanks(20, 17)
	match, err = userA.Challenge(userB)
	require.NoError(t, err)
	require.NotNil(t, match)
	requireMatchContains(t, match, userA, userB)

	userA, userB = getTestUsersWithRanks(30, 25)
	match, err = userA.Challenge(userB)
	require.NoError(t, err)
	require.NotNil(t, match)
	requireMatchContains(t, match, userA, userB)

	userA, userB = getTestUsersWithRanks(36, 28)
	match, err = userA.Challenge(userB)
	require.NoError(t, err)
	require.NotNil(t, match)
	requireMatchContains(t, match, userA, userB)
}

// test that a user cannot challenge/be challeneged with an active match
func TestChallengeInvalidActive(t *testing.T) {
	userA, userB := getTestUsersWithRanks(3, 4)
	var userC *User = models.NewUser("jdallard", "ncsu", "Jason", "Allard")
	userC.ID = 3
	userC.Rank = 5

	_, _ = userC.Challenge(userA)

	var match *Match
	var err error

	match, err = userB.Challenge(userC)

	require.Error(t, err)
	require.Nil(t, match)

	match, err = userA.Challenge(userB)

	require.Error(t, err)
	require.Nil(t, match)
}

// test that a user that is seeded too low cannot challenge a high-seeded player
func TestChallengeLowSeed(t *testing.T) {
	var userA *User
	var userB *User
	var match *Match
	var err error

	userA, userB = getTestUsersWithRanks(4, 1)
	match, err = userA.Challenge(userB)
	require.Error(t, err)
	require.Nil(t, match)
	require.Zero(t, len(userA.Matches))
	require.Zero(t, len(userB.Matches))

	userA, userB = getTestUsersWithRanks(15, 11)
	match, err = userA.Challenge(userB)
	require.Error(t, err)
	require.Nil(t, match)
	require.Zero(t, len(userA.Matches))
	require.Zero(t, len(userB.Matches))

	userA, userB = getTestUsersWithRanks(21, 15)
	match, err = userA.Challenge(userB)
	require.Error(t, err)
	require.Nil(t, match)
	require.Zero(t, len(userA.Matches))
	require.Zero(t, len(userB.Matches))

	userA, userB = getTestUsersWithRanks(35, 29)
	match, err = userA.Challenge(userB)
	require.Error(t, err)
	require.Nil(t, match)
	require.Zero(t, len(userA.Matches))
	require.Zero(t, len(userB.Matches))

	userA, userB = getTestUsersWithRanks(36, 27)
	match, err = userA.Challenge(userB)
	require.Error(t, err)
	require.Nil(t, match)
	require.Zero(t, len(userA.Matches))
	require.Zero(t, len(userB.Matches))
}

// test that ensures that a user can re-challenge after the cooldown expires
func TestChallengeTimeValid(t *testing.T) {
	const secondsInHour int = 3600

	userA, userB := getTestUsersWithRanks(3, 4)
	match, _ := userA.Challenge(userB)

	match.IsActive = false
	// player challenged player more than enough time ago
	match.CreatedAt = time.Unix(time.Now().Unix()-int64((models.SamePlayerCooldownHours+8)*secondsInHour), 0)

	match, err := userA.Challenge(userB)

	require.NoError(t, err)
	require.NotNil(t, match)
}

// test that ensures a player cannot re-challenge the same player while a cooldown is active
func TestChallengeTimeInValid(t *testing.T) {
	const secondsInHour int = 3600

	userA, userB := getTestUsersWithRanks(3, 4)
	match, _ := userA.Challenge(userB)

	match.IsActive = false
	// player challenged player more than enough time ago
	match.CreatedAt = time.Unix(time.Now().Unix()-int64((models.SamePlayerCooldownHours-8)*secondsInHour), 0)

	match, err := userA.Challenge(userB)

	require.Error(t, err)
	require.Nil(t, match)
}

// test that all possible scores can be encoded in 8 bits
func TestEncodeDecodeScore(t *testing.T) {
	scores := []int{1, 2, 3, 4, 5, 6}

	for _, i := range scores {
		for _, j := range scores {
			encoded := models.EncodeScore(i, j)
			id, jd := models.DecodeScore(encoded)
			require.Equal(t, i, id)
			require.Equal(t, j, jd)
		}
	}
}

func TestSubmitScoreValid(t *testing.T) {
	userA, userB := getTestUsersWithRanks(3, 4)
	match, _ := userA.Challenge(userB)

	err := match.SubmitScore(6, 4)

	require.NoError(t, err)
	require.False(t, match.IsActive)
	a, b := models.DecodeScore(match.Score)
	require.Equal(t, 6, a)
	require.Equal(t, 4, b)

}

func TestSubmitScoreInvalid(t *testing.T) {
	userA, userB := getTestUsersWithRanks(3, 4)
	match, _ := userA.Challenge(userB)

	err := match.SubmitScore(5, 5)
	require.Error(t, err)

	err = match.SubmitScore(6, -1)
	require.Error(t, err)

	err = match.SubmitScore(100, 6)
	require.Error(t, err)
}
