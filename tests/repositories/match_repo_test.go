package repositories_test

import (
	"ClubTennis/initializers"
	"ClubTennis/models"
	"ClubTennis/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type Match = models.Match
type User = models.User
type MatchRepository = repositories.MatchRepository

type MatchTestSuite struct {
	suite.Suite
	repo     *MatchRepository
	userRepo *repositories.UserRepository
	userA    *User
	userB    *User
}

// sets up before each test
func (suite *MatchTestSuite) SetupTest() {
	//drop the schema
	db := initializers.GetTestDatabase()
	if db == nil {
		panic("error in setup!")
	}
	db.Exec("DROP SCHEMA " + initializers.TestDBName + ";")

	//get the schema back. Burning my flash transistors.
	db = initializers.GetTestDatabase()

	err := db.AutoMigrate(&Match{}, &User{})
	if err != nil {
		panic(err)
	}

	suite.userRepo = repositories.NewUserRepository(db)
	suite.repo = repositories.NewMatchRepository(db)

	suite.userA, _ = models.NewUser("bdoller4", "ncsu", "bowie", "doliver", "bdoller4@ncsu.edu", models.MENS_LADDER)
	suite.userA.Rank = 4
	suite.userB, _ = models.NewUser("qbingus5", "ncsu", "quevin", "bingus", "qbingus5@ncsu.edu", models.MENS_LADDER)
	suite.userB.Rank = 5

	if err = suite.userRepo.SubmitUser(suite.userA); err != nil {
		panic("err in setup!")
	}
	if err = suite.userRepo.SubmitUser(suite.userB); err != nil {
		panic("err in setup!")
	}

	if suite.repo == nil || suite.userA == nil || suite.userB == nil {
		panic("error in setup!")
	}

}

// neccessary for 'go test' to call all suite tests
func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(MatchTestSuite))
}

func (suite *MatchTestSuite) TestGetMatch() {
	match, err := suite.userA.Challenge(suite.userB)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(match)

	err = suite.repo.SubmitMatch(match)
	suite.Require().NoError(err)

	fetched, err := suite.repo.FindByID(match.ID)
	suite.Require().NoError(err)
	suite.Require().NotNil(fetched)

	suite.Require().Equal(match.Challenger().UnityID, fetched.Challenger().UnityID)
	suite.Require().Equal(match.Challenged().UnityID, fetched.Challenged().UnityID)

}

func (suite *MatchTestSuite) TestSubmitScoreDB() {
	match, _ := suite.userA.Challenge(suite.userB)

	err := suite.repo.SubmitMatch(match)
	suite.Assert().NoError(err)

	fetched, err := suite.repo.FindByID(match.ID)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(fetched)

	err = fetched.SubmitScore(6, 4)
	suite.Assert().NoError(err)

	suite.repo.SaveMatch(fetched)

	f2, err := suite.repo.FindByID(fetched.ID)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(f2)

	a, b := models.DecodeScore(f2.Score)

	suite.Require().Equal(6, a)
	suite.Require().Equal(4, b)
	suite.Require().False(f2.IsActive)

	userA, _ := suite.userRepo.FindByID(suite.userA.ID)
	userB, _ := suite.userRepo.FindByID(suite.userB.ID)

	suite.Require().Equal(1, userA.Wins)
	suite.Assert().Equal(0, userA.Losses)
	suite.Require().Equal(0, userB.Wins)
	suite.Assert().Equal(1, userB.Losses)
}

func (suite *MatchTestSuite) TestFindMatchByUsers() {
	match, _ := suite.userB.Challenge(suite.userA)
	suite.repo.SubmitMatch(match)

	fetched, err := suite.repo.FindByChallengerChallenged(suite.userB.ID, suite.userA.ID)
	suite.Require().NoError(err)
	suite.Require().NotNil(fetched)
	suite.Assert().Len(fetched, 1)
	suite.Require().True(fetched[0].IsActive)
	//fudge the time so it looks like it was made a while ago
	fetched[0].CreatedAt = time.Unix(time.Now().Unix()-int64((models.SamePlayerCooldownHours+300)*3600), 0)
	fetched[0].SubmitScore(6, 3)
	suite.repo.SaveMatch(&fetched[0])

	suite.userA, _ = suite.userRepo.FindByUnityID(suite.userA.UnityID)
	suite.userB, _ = suite.userRepo.FindByUnityID(suite.userB.UnityID)

	match, err = suite.userA.Challenge(suite.userB)
	suite.Require().NoError(err)
	suite.Require().NotNil(match)

	suite.repo.SubmitMatch(match)

	fetched, err = suite.repo.FindByPlayerIDs(suite.userA.ID)
	suite.Require().Len(fetched, 2)
	suite.Require().NoError(err)

	fetched, err = suite.repo.FindByPlayerIDs(suite.userB.ID)
	suite.Require().Len(fetched, 2)
	suite.Require().NoError(err)

	fetched, err = suite.repo.FindByPlayerIDs(suite.userA.ID, suite.userB.ID)
	suite.Require().Len(fetched, 2)
	suite.Require().NoError(err)

}
