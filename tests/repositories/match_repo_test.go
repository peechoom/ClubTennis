package repositories_test

import (
	"ClubTennis/initializers"
	"ClubTennis/models"
	"ClubTennis/repositories"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Match = models.Match
type User = models.User
type MatchRepository = repositories.MatchRepository

type RepoTestSuite struct {
	suite.Suite
	repo     *MatchRepository
	userRepo *repositories.UserRepository
	userA    *User
	userB    *User
}

// sets up before each test
func (suite *RepoTestSuite) SetupTest() {
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

	suite.userA, _ = models.NewUser("bdoller4", "ncsu", "bowie", "doliver", "bdoller4@ncsu.edu")
	suite.userA.Rank = 4
	suite.userB, _ = models.NewUser("qbingus5", "ncsu", "quevin", "bingus", "qbingus5@ncsu.edu")
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
	suite.Run(t, new(RepoTestSuite))
}

func (suite *RepoTestSuite) TestGetMatch() {
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
