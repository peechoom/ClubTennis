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
	repo  *MatchRepository
	userA *User
	userB *User
}

// sets up before each test
func (suite *RepoTestSuite) SetupTest() {
	//drop the schema
	db := initializers.GetTestDatabase()
	if db == nil {
		panic("error in setup!")
	}
	db.Exec("DROP SCHEMA " + initializers.TestDBName + ";")

	//get the schema back
	db = initializers.GetTestDatabase()

	err := db.AutoMigrate(&Match{})
	if err != nil {
		panic(err)
	}
	suite.repo = repositories.NewMatchRepository(db)
	suite.userA, _ = models.NewUser("bdoller4", "ncsu", "bowie", "doliver", "bdoller4@ncsu.edu")
	suite.userA.Rank = 4
	suite.userB, _ = models.NewUser("qbingus5", "ncsu", "quevin", "bingus", "qbingus5@ncsu.edu")
	suite.userB.Rank = 5

	//TODO save these two

	if suite.repo == nil || suite.userA == nil || suite.userB == nil {
		panic("error in setup!")
	}

}

// neccessary for 'go test' to call all suite tests
func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}

func (suite *RepoTestSuite) TestGetMatch() {

}
