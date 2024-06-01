package services_test

import (
	"ClubTennis/initializers"
	"ClubTennis/models"
	"ClubTennis/services"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	s     *services.UserService
	userA *models.User
	userB *models.User
	userC *models.User
	userD *models.User
	userE *models.User
	userF *models.User
	userG *models.User
}

// sets up before each test
func (suite *UserServiceTestSuite) SetupTest() {
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

	suite.s = services.NewUserService(db)
	suite.userA, _ = models.NewUser("shboil4", "ncsu", "Sam", "Boiland", "shboil4@ncsu.edu", models.MENS_LADDER)
	suite.userB, _ = models.NewUser("jbeno5", "ncsu", "James", "Benolli", "jbeno5@ncsu.edu", models.MENS_LADDER)
	suite.userC, _ = models.NewUser("pdiddy4", "ncsu", "Puff", "Daddy", "pdiddy@ncsu.edu", models.MENS_LADDER)
	suite.userD, _ = models.NewUser("jobitch2", "ncsu", "Joel", "Embitch", "jobitch@ncsu.edu", models.MENS_LADDER)
	suite.userE, _ = models.NewUser("myprince2", "ncsu", "Lebron", "James", "myprince2@ncsu.edu", models.MENS_LADDER)
	suite.userF, _ = models.NewUser("alhoot4", "ncsu", "Alison", "Hoot", "alhoot4@ncsu.edu", models.WOMENS_LADDER)
	suite.userG, _ = models.NewUser("lray1", "ncsu", "Lana", "DelRay", "lray1@ncsu.edu", models.WOMENS_LADDER)

	suite.userF.Rank = 1
	suite.userG.Rank = 2
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
	suite.userF.Matches = make([]*models.Match, 0)
	suite.userG.Matches = make([]*models.Match, 0)

}

// neccessary for 'go test' to call all suite tests
func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) TestServiceSaveFind() {
	suite.s.Save(suite.userA)
	var u []models.User
	u, err := suite.s.FindAll()
	suite.Require().NoError(err)
	suite.Assert().Len(u, 1)
	suite.Require().Equal(suite.userA.UnityID, u[0].UnityID)

	suite.s.Save(suite.userB, suite.userC, suite.userD, suite.userE)

	u, err = suite.s.FindByRankRange(models.MENS_LADDER, 1, 5)
	suite.Require().NoError(err)
	suite.Assert().Len(u, 5)

}

func (suite *UserServiceTestSuite) TestServiceFindByUnityID() {
	suite.s.Save(suite.userA)
	suite.s.Save(suite.userB)
	suite.s.Save(suite.userC)
	suite.s.Save(suite.userD)
	suite.s.Save(suite.userE)

	var search *models.User
	var err error

	search, err = suite.s.FindByUnityID("shboil4")
	suite.Require().NoError(err)
	suite.Require().Equal(suite.userA, search)

	search, err = suite.s.FindByUnityID("farthead1")
	suite.Assert().NoError(err)
	suite.Require().Nil(search)
}

func (suite *UserServiceTestSuite) TestLadderAlgo() {
	suite.s.Save(suite.userA, suite.userB, suite.userC, suite.userD, suite.userE, suite.userF, suite.userG)
	suite.userA, _ = suite.s.FindByUnityID(suite.userA.UnityID)
	suite.userB, _ = suite.s.FindByUnityID(suite.userB.UnityID)
	suite.userC, _ = suite.s.FindByUnityID(suite.userC.UnityID)
	suite.userD, _ = suite.s.FindByUnityID(suite.userD.UnityID)
	suite.userE, _ = suite.s.FindByUnityID(suite.userE.UnityID)
	suite.userF, _ = suite.s.FindByUnityID(suite.userF.UnityID)
	suite.userG, _ = suite.s.FindByUnityID(suite.userG.UnityID)

	suite.s.AdjustLadder(suite.userC, suite.userA)

	suite.userA, _ = suite.s.FindByUnityID(suite.userA.UnityID)
	suite.userB, _ = suite.s.FindByUnityID(suite.userB.UnityID)
	suite.userC, _ = suite.s.FindByUnityID(suite.userC.UnityID)
	suite.userD, _ = suite.s.FindByUnityID(suite.userD.UnityID)
	suite.userE, _ = suite.s.FindByUnityID(suite.userE.UnityID)
	suite.userF, _ = suite.s.FindByUnityID(suite.userF.UnityID)
	suite.userG, _ = suite.s.FindByUnityID(suite.userG.UnityID)

	suite.Assert().Equal(uint(1), suite.userC.Rank)
	suite.Assert().Equal(uint(2), suite.userA.Rank)
	suite.Assert().Equal(uint(3), suite.userB.Rank)
	suite.Assert().Equal(uint(4), suite.userD.Rank)
	suite.Assert().Equal(uint(5), suite.userE.Rank)

	suite.Assert().Equal(uint(1), suite.userF.Rank)
	suite.Assert().Equal(uint(2), suite.userG.Rank)

	suite.s.AdjustLadder(suite.userC, suite.userA)

	suite.userA, _ = suite.s.FindByUnityID(suite.userA.UnityID)
	suite.userB, _ = suite.s.FindByUnityID(suite.userB.UnityID)
	suite.userC, _ = suite.s.FindByUnityID(suite.userC.UnityID)
	suite.userD, _ = suite.s.FindByUnityID(suite.userD.UnityID)
	suite.userE, _ = suite.s.FindByUnityID(suite.userE.UnityID)
	suite.userF, _ = suite.s.FindByUnityID(suite.userF.UnityID)
	suite.userG, _ = suite.s.FindByUnityID(suite.userG.UnityID)

	suite.Assert().Equal(uint(1), suite.userC.Rank)
	suite.Assert().Equal(uint(2), suite.userA.Rank)
	suite.Assert().Equal(uint(3), suite.userB.Rank)
	suite.Assert().Equal(uint(4), suite.userD.Rank)
	suite.Assert().Equal(uint(5), suite.userE.Rank)

	suite.Assert().Equal(uint(1), suite.userF.Rank)
	suite.Assert().Equal(uint(2), suite.userG.Rank)

	suite.s.AdjustLadder(suite.userA, suite.userC)

	suite.userA, _ = suite.s.FindByUnityID(suite.userA.UnityID)
	suite.userB, _ = suite.s.FindByUnityID(suite.userB.UnityID)
	suite.userC, _ = suite.s.FindByUnityID(suite.userC.UnityID)
	suite.userD, _ = suite.s.FindByUnityID(suite.userD.UnityID)
	suite.userE, _ = suite.s.FindByUnityID(suite.userE.UnityID)
	suite.userF, _ = suite.s.FindByUnityID(suite.userF.UnityID)
	suite.userG, _ = suite.s.FindByUnityID(suite.userG.UnityID)

	suite.Assert().Equal(uint(1), suite.userA.Rank)
	suite.Assert().Equal(uint(2), suite.userC.Rank)
	suite.Assert().Equal(uint(3), suite.userB.Rank)
	suite.Assert().Equal(uint(4), suite.userD.Rank)
	suite.Assert().Equal(uint(5), suite.userE.Rank)

	suite.Assert().Equal(uint(1), suite.userF.Rank)
	suite.Assert().Equal(uint(2), suite.userG.Rank)

	suite.s.AdjustLadder(suite.userE, suite.userD)

	suite.userA, _ = suite.s.FindByUnityID(suite.userA.UnityID)
	suite.userB, _ = suite.s.FindByUnityID(suite.userB.UnityID)
	suite.userC, _ = suite.s.FindByUnityID(suite.userC.UnityID)
	suite.userD, _ = suite.s.FindByUnityID(suite.userD.UnityID)
	suite.userE, _ = suite.s.FindByUnityID(suite.userE.UnityID)
	suite.userF, _ = suite.s.FindByUnityID(suite.userF.UnityID)
	suite.userG, _ = suite.s.FindByUnityID(suite.userG.UnityID)

	suite.Assert().Equal(uint(1), suite.userA.Rank)
	suite.Assert().Equal(uint(2), suite.userC.Rank)
	suite.Assert().Equal(uint(3), suite.userB.Rank)
	suite.Assert().Equal(uint(4), suite.userE.Rank)
	suite.Assert().Equal(uint(5), suite.userD.Rank)

	suite.Assert().Equal(uint(1), suite.userF.Rank)
	suite.Assert().Equal(uint(2), suite.userG.Rank)

	suite.s.AdjustLadder(suite.userD, suite.userA)

	suite.userA, _ = suite.s.FindByUnityID(suite.userA.UnityID)
	suite.userB, _ = suite.s.FindByUnityID(suite.userB.UnityID)
	suite.userC, _ = suite.s.FindByUnityID(suite.userC.UnityID)
	suite.userD, _ = suite.s.FindByUnityID(suite.userD.UnityID)
	suite.userE, _ = suite.s.FindByUnityID(suite.userE.UnityID)
	suite.userF, _ = suite.s.FindByUnityID(suite.userF.UnityID)
	suite.userG, _ = suite.s.FindByUnityID(suite.userG.UnityID)

	suite.Assert().Equal(uint(1), suite.userD.Rank)
	suite.Assert().Equal(uint(2), suite.userA.Rank)
	suite.Assert().Equal(uint(3), suite.userC.Rank)
	suite.Assert().Equal(uint(4), suite.userB.Rank)
	suite.Assert().Equal(uint(5), suite.userE.Rank)

	suite.Assert().Equal(uint(1), suite.userF.Rank)
	suite.Assert().Equal(uint(2), suite.userG.Rank)

	suite.s.AdjustLadder(suite.userG, suite.userF)

	suite.userA, _ = suite.s.FindByUnityID(suite.userA.UnityID)
	suite.userB, _ = suite.s.FindByUnityID(suite.userB.UnityID)
	suite.userC, _ = suite.s.FindByUnityID(suite.userC.UnityID)
	suite.userD, _ = suite.s.FindByUnityID(suite.userD.UnityID)
	suite.userE, _ = suite.s.FindByUnityID(suite.userE.UnityID)
	suite.userF, _ = suite.s.FindByUnityID(suite.userF.UnityID)
	suite.userG, _ = suite.s.FindByUnityID(suite.userG.UnityID)

	suite.Assert().Equal(uint(1), suite.userD.Rank)
	suite.Assert().Equal(uint(2), suite.userA.Rank)
	suite.Assert().Equal(uint(3), suite.userC.Rank)
	suite.Assert().Equal(uint(4), suite.userB.Rank)
	suite.Assert().Equal(uint(5), suite.userE.Rank)

	suite.Assert().Equal(uint(2), suite.userF.Rank)
	suite.Assert().Equal(uint(1), suite.userG.Rank)

}

func (suite *UserServiceTestSuite) TestDelete() {
	suite.s.Save(suite.userA, suite.userB, suite.userC, suite.userD, suite.userE)

	err := suite.s.DeleteByUnityID(suite.userB.UnityID)
	suite.Require().NoError(err)

	A, err := suite.s.FindByUnityID(suite.userA.UnityID)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(A)
	B, err := suite.s.FindByUnityID(suite.userB.UnityID)
	suite.Assert().NoError(err)
	suite.Assert().Nil(B)
	C, err := suite.s.FindByUnityID(suite.userC.UnityID)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(C)
	E, err := suite.s.FindByUnityID(suite.userE.UnityID)
	suite.Assert().NoError(err)
	suite.Assert().NotNil(E)

	suite.Require().Equal(uint(1), A.Rank)
	suite.Require().Equal(uint(2), C.Rank)
	suite.Require().Equal(uint(4), E.Rank)

}
