package repositories_test

import (
	"ClubTennis/initializers"
	"ClubTennis/models"
	"ClubTennis/repositories"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	repo  *repositories.UserRepository
	mr    *repositories.MatchRepository
	userA *User
	userB *User
}

// sets up before each test
func (suite *UserTestSuite) SetupTest() {
	//drop the schema
	db := initializers.GetTestDatabase()
	if db == nil {
		panic("error in setup!")
	}
	db.Exec("DROP SCHEMA " + initializers.TestDBName + ";")

	//get the schema back
	db = initializers.GetTestDatabase()

	err := db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

	suite.repo = repositories.NewUserRepository(db)
	suite.mr = repositories.NewMatchRepository(db)
	suite.userA, _ = models.NewUser("bdoller4", "ncsu", "bowie", "doliver", "bdoller4@ncsu.edu")
	suite.userB, _ = models.NewUser("qbingus5", "ncsu", "quevin", "bingus", "qbingus5@ncsu.edu")
	suite.userA.Matches = make([]*models.Match, 0)
	suite.userB.Matches = make([]*models.Match, 0)

	if suite.repo == nil || suite.userA == nil || suite.userB == nil {
		panic("error in setup!")
	}

}

// neccessary for 'go test' to call all suite tests
func TestUserRepoSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) TestGetUser() {
	err := suite.repo.SubmitUser(suite.userA)

	suite.Require().NoError(err)
	suite.Require().NotZero(suite.userA.ID)
	var id uint = suite.userA.ID
	var fetchedUser *User

	fetchedUser, err = suite.repo.FindByID(id)

	suite.Require().NoError(err)
	suite.Require().Equal(suite.userA, fetchedUser)
}

func (suite *UserTestSuite) TestSaveUser() {
	userA := suite.userA
	err := suite.repo.SubmitUser(userA)
	suite.Assert().NoError(err)

	userA.Affiliation = "skema"
	var createdUser *User
	var updatedUser *User

	createdUser, _ = suite.repo.FindByID(userA.ID)
	suite.Assert().NotEqual(userA, createdUser)

	err = suite.repo.SaveUser(userA)
	suite.Assert().NoError(err)
	updatedUser, err = suite.repo.FindByID(userA.ID)
	suite.Assert().NoError(err)

	suite.Require().Equal(userA, updatedUser)
}

func (suite *UserTestSuite) TestUserBatchSave() {
	var s []models.User
	s = append(s, *suite.userA)
	s = append(s, *suite.userB)

	err := suite.repo.SubmitUsers(s)

	suite.Require().NoError(err)
	a, err := suite.repo.FindByID(s[0].ID)
	suite.Require().NoError(err)
	suite.Require().Equal(s[0], *a)

	b, err := suite.repo.FindByID(s[1].ID)
	suite.Require().NoError(err)
	suite.Require().Equal(s[1], *b)
}

func (suite *UserTestSuite) TestGetByRanking() {
	suite.userA.Rank = 3
	suite.userB.Rank = 4

	suite.repo.SaveUser(suite.userA)
	suite.repo.SaveUser(suite.userB)

	a, err := suite.repo.FindByRank(3)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.userA.UnityID, a.UnityID)
	b, err := suite.repo.FindByRank(4)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.userB.UnityID, b.UnityID)

	var u []models.User
	u, err = suite.repo.FindByRankRange(3, 4)

	sort.Slice(u, func(i, j int) bool {
		return u[i].Rank < u[j].Rank
	})

	suite.Require().NoError(err)
	// is a pointer... mumble mumble
	suite.Require().Equal(u[0].UnityID, (suite.userA).UnityID)
	suite.Require().Equal(u[1].UnityID, (suite.userB).UnityID)

}

func (suite *UserTestSuite) TestSubmitUser() {
	err := suite.repo.SubmitUser(suite.userA)
	suite.Require().NoError(err)
	ua, _ := suite.repo.FindByUnityID(suite.userA.UnityID)
	ua.Rank = 1
	suite.repo.SaveUser(ua)

	err = suite.repo.SubmitUser(suite.userB)
	suite.Require().NoError(err)
	ub, _ := suite.repo.FindByUnityID(suite.userB.UnityID)
	suite.Require().Equal(uint(2), ub.Rank)

	uc, _ := models.NewOfficer("kwest4", "ncsu", "Kanye", "West", "kwest4@ncsu.edu")
	err = suite.repo.SubmitUser(uc)
	suite.Require().Equal(uint(3), uc.Rank)

}

func (suite *UserTestSuite) TestDeleteUser() {
	err := suite.repo.SubmitUser(suite.userA)
	suite.Assert().NoError(err)

	err = suite.repo.DeleteByID(suite.userA.ID)
	suite.Require().NoError(err)

	everything, _ := suite.repo.FindAll()
	suite.Require().Zero(len(everything))
	suite.userA.ID = 0
	suite.Assert().NoError(suite.repo.SubmitUser(suite.userA))
	suite.Assert().NoError(suite.repo.SubmitUser(suite.userB))

	m, _ := suite.userA.Challenge(suite.userB)
	suite.mr.SubmitMatch(m)
	suite.repo.SaveUser(suite.userA)
	suite.repo.SaveUser(suite.userB)

	err = suite.repo.DeleteByID(suite.userA.ID)
	suite.Require().NoError(err)
	everything, _ = suite.repo.FindAll()
	suite.Require().Len(everything, 1)

	b, _ := suite.repo.FindByID(suite.userB.ID)
	suite.Require().Len(b.Matches, 1)
	suite.Require().Nil(b.Matches[0].Challenger())

	a, _ := suite.repo.FindByID(suite.userA.ID)
	suite.Require().Nil(a)
}
