package services_test

import (
	"ClubTennis/initializers"
	"ClubTennis/models"
	"ClubTennis/services"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MatchServiceTestSuite struct {
	suite.Suite
	s     *services.MatchService
	us    *services.UserService
	userA *models.User
	userB *models.User
	userC *models.User
	userD *models.User
	userE *models.User
}

func (suite *MatchServiceTestSuite) SetupTest() {
	db := initializers.GetTestDatabase()
	if db == nil {
		panic("error in setup!")
	}
	db.Exec("DROP SCHEMA " + initializers.TestDBName + ";")

	db = initializers.GetTestDatabase()
	if db == nil {
		panic("error in setup!")
	}
	err := db.AutoMigrate(models.User{}, models.Match{})
	if err != nil {
		panic(err)
	}

	suite.s = services.NewMatchService(db)
	suite.us = services.NewUserService(db)

	suite.userA, _ = models.NewUser("shboil4", "ncsu", "Sam", "Boiland", "shboil4@ncsu.edu")
	suite.userB, _ = models.NewUser("jbeno5", "ncsu", "James", "Benolli", "jbeno5@ncsu.edu")
	suite.userC, _ = models.NewUser("pdiddy4", "ncsu", "Puff", "Daddy", "pdiddy@ncsu.edu")
	suite.userD, _ = models.NewUser("jobitch2", "ncsu", "Joel", "Embitch", "jobitch@ncsu.edu")
	suite.userE, _ = models.NewUser("myprince2", "ncsu", "Lebron", "James", "myprince2@ncsu.edu")

	suite.userA.Rank = 1
	suite.userB.Rank = 2
	suite.userC.Rank = 3
	suite.userD.Rank = 4
	suite.userE.Rank = 5

	suite.userA.Matches = make([]*models.Match, 0)
	suite.userB.Matches = make([]*models.Match, 0)
	suite.userC.Matches = make([]*models.Match, 0)
	suite.userD.Matches = make([]*models.Match, 0)
	suite.userE.Matches = make([]*models.Match, 0)

	suite.us.Save(suite.userA)
	suite.us.Save(suite.userB)
	suite.us.Save(suite.userC)
	suite.us.Save(suite.userD)
	suite.us.Save(suite.userE)

}

func TestMatchServiceSuite(t *testing.T) {
	suite.Run(t, new(MatchServiceTestSuite))
}

func (suite *MatchServiceTestSuite) TestSaveNewMatch() {
	userA, _ := suite.us.FindByID(suite.userA.ID)
	userB, _ := suite.us.FindByID(suite.userB.ID)

	match, err := userA.Challenge(userB)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(match)

	err = suite.s.Save(match)
	suite.Require().NoError(err)

	fetchedMatch, err := suite.s.FindByPlayerID(userA.ID, userB.ID)
	suite.Require().NoError(err)
	suite.Require().Len(fetchedMatch, 1)
	suite.Require().Equal(userA.ID, fetchedMatch[0].ChallengerID)
	suite.Require().Equal(userB.ID, fetchedMatch[0].ChallengedID)

	userA, _ = suite.us.FindByID(userA.ID)
	userB, _ = suite.us.FindByID(userB.ID)

	suite.Require().Equal(userA.Matches[0].ChallengerID, fetchedMatch[0].ChallengerID)
	suite.Require().Equal(userA.Matches[0].ChallengedID, fetchedMatch[0].ChallengedID)
	suite.Require().Equal(userB.Matches[0].ChallengerID, fetchedMatch[0].ChallengerID)
	suite.Require().Equal(userB.Matches[0].ChallengedID, fetchedMatch[0].ChallengedID)

}

func (suite *MatchServiceTestSuite) TestSubmitMatchScore() {
	userD, _ := suite.us.FindByID(suite.userD.ID)
	userB, _ := suite.us.FindByID(suite.userB.ID)

	match, err := userD.Challenge(userB)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(match)
	err = suite.s.Save(match)
	suite.Assert().NoError(err)

	fetchedMatch, err := suite.s.FindByPlayerID(userD.ID, userB.ID)
	suite.Assert().NoError(err)
	suite.Require().Len(fetchedMatch, 1)

	err = fetchedMatch[0].SubmitScore(6, 2)
	suite.Assert().NoError(err)
	suite.s.Save(&fetchedMatch[0])

	fetchedMatch, err = suite.s.FindByPlayerID(userD.ID, userB.ID)
	suite.Assert().NoError(err)
	suite.Require().Len(fetchedMatch, 1)

	suite.Require().False(fetchedMatch[0].IsActive)
	suite.Require().Equal(models.EncodeScore(6, 2), fetchedMatch[0].Score)

	userA, _ := suite.us.FindByID(suite.userA.ID)
	userB, _ = suite.us.FindByID(suite.userB.ID)
	userC, _ := suite.us.FindByID(suite.userC.ID)
	userD, _ = suite.us.FindByID(suite.userD.ID)
	userE, _ := suite.us.FindByID(suite.userE.ID)

	suite.Require().Equal(uint(1), userA.Rank)
	suite.Require().Equal(uint(2), userD.Rank)
	suite.Require().Equal(uint(3), userB.Rank)
	suite.Require().Equal(uint(4), userC.Rank)
	suite.Require().Equal(uint(5), userE.Rank)
}
