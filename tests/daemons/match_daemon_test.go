package daemons_test

import (
	"ClubTennis/config"
	"ClubTennis/daemons"
	"ClubTennis/initializers"
	"ClubTennis/models"
	"ClubTennis/services"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type matchDaemonSuite struct {
	suite.Suite
	ms    *services.MatchService
	es    *services.EmailService
	us    *services.UserService
	userA *models.User
	userB *models.User
}

func (suite *matchDaemonSuite) SetupTest() {
	config.LoadConfig("/home/alec/go/src/ClubTennis/config/.env")
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
	suite.us = services.NewUserService(db)
	suite.ms = services.NewMatchService(db)
	suite.es = services.NewEmailService("/home/alec/go/src/ClubTennis/templates", "a", "a")

	suite.userA, _ = models.NewUser("pjdiddy4", "ncsu", "Puff", "Diddy", "null@gmail.com", models.MENS_LADDER)
	suite.userB, _ = models.NewUser("tgrizzly", "ncsu", "Tee", "Grizzly", "nil@gmail.com", models.MENS_LADDER)
	suite.userA.Rank = 1
	suite.userB.Rank = 2
	suite.userA.Matches = make([]*models.Match, 0)
	suite.userB.Matches = make([]*models.Match, 0)
	suite.us.Save(suite.userA)
	suite.us.Save(suite.userB)

}

func TestMatchDaemons(t *testing.T) {
	suite.Run(t, new(matchDaemonSuite))
}

func (suite *matchDaemonSuite) TestWarningDaemon() {
	var err error
	suite.userA, err = suite.us.FindByUnityID(suite.userA.UnityID)
	suite.Require().NoError(err)
	suite.userB, err = suite.us.FindByUnityID(suite.userB.UnityID)
	suite.Require().NoError(err)

	m, _ := suite.userB.Challenge(suite.userA)
	suite.Require().NotNil(m)
	m.SubmittedAt = time.Now().Add(-(time.Hour*24*8 + time.Hour))

	suite.ms.Save(m)

	found, err := suite.ms.FindByChallengerID(suite.userB.ID)
	m = &found[0]
	suite.Require().NoError(err)
	suite.Require().True(m.SubmittedAt.Before(time.Now().Add(-(time.Hour * 24 * 8))))
	suite.Require().True(m.SubmittedAt.After(time.Now().Add(-(time.Hour * 24 * 9))))

	w, err := suite.ms.FindByNearlyExpired()
	suite.Require().NoError(err)
	suite.Require().Len(w, 1)
	suite.Require().Equal(m.ChallengerID, w[0].ChallengerID)
	suite.Require().Equal(m.ChallengedID, w[0].ChallengedID)

	go daemons.MatchWarningDaemon(time.Second, false, suite.ms, suite.es)

	time.Sleep(time.Second * 2)
	m, err = suite.ms.FindByID(m.ID)
	suite.Require().NoError(err)
	suite.Require().True(m.LateNotifSent)

}

func (suite *matchDaemonSuite) TestForfeitDaemon() {
	var err error
	suite.userA, err = suite.us.FindByUnityID(suite.userA.UnityID)
	suite.Require().NoError(err)
	suite.userB, err = suite.us.FindByUnityID(suite.userB.UnityID)
	suite.Require().NoError(err)

	m, _ := suite.userB.Challenge(suite.userA)
	suite.Require().NotNil(m)
	m.SubmittedAt = time.Now().Add(-(time.Hour * 24 * 9))

	suite.ms.Save(m)

	go daemons.MatchExpiredDaemon(time.Second, false, suite.ms, suite.es)

	time.Sleep(time.Second * 2)
	m, err = suite.ms.FindByID(m.ID)
	suite.Require().NoError(err)
	suite.Require().Equal(models.EncodeScore(6, 0), m.Score)
	suite.Require().False(m.IsActive)
}

func (suite *matchDaemonSuite) TestDeletionDaemon() {
	var err error
	suite.userA, err = suite.us.FindByUnityID(suite.userA.UnityID)
	suite.Require().NoError(err)
	suite.userB, err = suite.us.FindByUnityID(suite.userB.UnityID)
	suite.Require().NoError(err)

	m, _ := suite.userB.Challenge(suite.userA)
	suite.Require().NotNil(m)
	m.SubmitScore(6, 4)
	m.SubmittedAt = time.Now().Add(-(time.Hour * 24 * 60))

	suite.ms.Save(m)

	go daemons.MatchDeletionDaemon(time.Second, false, suite.ms)

	time.Sleep(time.Second * 2)

	m, err = suite.ms.FindByID(m.ID)
	suite.Require().Nil(m)
	suite.Require().Error(err)
}
