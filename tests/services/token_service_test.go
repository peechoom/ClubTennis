package services_test

import (
	"ClubTennis/config"
	"ClubTennis/models"
	"ClubTennis/repositories"
	"ClubTennis/services"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/mazen160/go-random"
	"github.com/stretchr/testify/suite"
)

type TokenTestSuite struct {
	suite.Suite
	ts    *services.TokenService
	userA *models.User
	userB *models.User
}

func (suite *TokenTestSuite) SetupTest() {
	err := config.LoadConfig("/home/alec/go/src/ClubTennis/config/.env")
	if err != nil {
		panic(err.Error())
	}

	suite.ts = services.DefaultTokenService(repositories.NewTokenRepository())
	suite.userA, _ = models.NewOfficer("kwest4", "ncsu", "Kanye", "West", "kwest4@ncsu.edu")
	suite.userA.ID = 20
	suite.userB, _ = models.NewUser("cwatts3", "unc", "Chris", "Watts", "cwatts3@unc.edu") //yes he went to UNC
	suite.userB.ID = 87

	suite.userA.Matches = make([]*models.Match, 0)
	suite.userB.Matches = make([]*models.Match, 0)

}

func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}

func (suite *TokenTestSuite) TestGetNewTokenPair() {
	tp, err := suite.ts.GetNewTokenPair(suite.userA.ID, "")
	suite.Require().NoError(err)
	suite.Assert().NotNil(tp)
	userID, err := suite.ts.ValidateIDToken(tp.IDToken.SS)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.userA.ID, userID)

	refreshToken, err := suite.ts.ValidateRefreshToken(tp.RefreshToken.SS)
	suite.Require().NoError(err)
	suite.Require().Equal(*refreshToken, tp.RefreshToken)

	tp, err = suite.ts.GetNewTokenPair(suite.userA.ID, refreshToken.ID.String())
	suite.Require().NoError(err)
	suite.Assert().NotNil(tp)

	userID, err = suite.ts.ValidateIDToken(tp.IDToken.SS)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.userA.ID, userID)

	refreshToken, err = suite.ts.ValidateRefreshToken(tp.RefreshToken.SS)
	suite.Require().NoError(err)
	suite.Require().Equal(*refreshToken, tp.RefreshToken)
}

func (suite *TokenTestSuite) TestExpiredToken() {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	sym, _ := random.Bytes(64)

	suite.ts = services.NewTokenService(repositories.NewTokenRepository(),
		priv, priv.Public().(*rsa.PublicKey), sym, 1, 2)

	tp, err := suite.ts.GetNewTokenPair(suite.userA.ID, "")
	suite.Require().NoError(err)
	suite.Assert().NotNil(tp)
	userID, err := suite.ts.ValidateIDToken(tp.IDToken.SS)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.userA.ID, userID)

	time.Sleep(time.Second)
	userID, err = suite.ts.ValidateIDToken(tp.IDToken.SS)
	suite.Require().Error(err)
	suite.Require().Zero(userID)

	refreshToken, err := suite.ts.ValidateRefreshToken(tp.RefreshToken.SS)
	suite.Require().NoError(err)
	suite.Require().Equal(*refreshToken, tp.RefreshToken)

	time.Sleep(time.Second)
	//this panics instead of just returning an error when reassigning tp ...weird
	badPair, err := suite.ts.GetNewTokenPair(suite.userA.ID, refreshToken.ID.String())
	suite.Require().Error(err)
	suite.Require().Nil(badPair)

	refreshToken, err = suite.ts.ValidateRefreshToken(tp.RefreshToken.SS)
	suite.Require().Error(err)
	suite.Require().Nil(refreshToken)

}
