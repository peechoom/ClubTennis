package services_test

import (
	"ClubTennis/config"
	"ClubTennis/initializers"
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
	c     *config.Config
	userA *models.User
	userB *models.User
}

func (suite *TokenTestSuite) SetupTest() {
	var err error
	suite.c, err = config.LoadConfig("/home/alec/go/src/ClubTennis/config/config.json")
	if err != nil || suite.c == nil {
		panic(err.Error())
	}
	client := initializers.GetTestClient(suite.c)
	if client == nil {
		panic("couldnt get client")
	}
	client.FlushAll()
	suite.ts = services.DefaultTokenService(repositories.NewTokenRepository(client))
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
	tp, err := suite.ts.GetNewTokenPair(suite.userA, "")
	suite.Require().NoError(err)
	suite.Assert().NotNil(tp)
	userID, err := suite.ts.ValidateIDToken(tp.IDToken.SS)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.userA.ID, userID)

	refreshToken, err := suite.ts.ValidateRefreshToken(tp.RefreshToken.SS)
	suite.Require().NoError(err)
	suite.Require().Equal(*refreshToken, tp.RefreshToken)

	tp, err = suite.ts.GetNewTokenPair(suite.userA, refreshToken.ID.String())
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

	suite.ts = services.NewTokenService(repositories.NewTokenRepository(initializers.GetTestClient(suite.c)),
		priv, priv.Public().(*rsa.PublicKey), sym, 1, 2)

	tp, err := suite.ts.GetNewTokenPair(suite.userA, "")
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
	badPair, err := suite.ts.GetNewTokenPair(suite.userA, refreshToken.ID.String())
	suite.Require().Error(err)
	suite.Require().Nil(badPair)

	refreshToken, err = suite.ts.ValidateRefreshToken(tp.RefreshToken.SS)
	suite.Require().Error(err)
	suite.Require().Nil(refreshToken)

}

func (suite *TokenTestSuite) TestDeleteAllTokens() {
	tp1, _ := suite.ts.GetNewTokenPair(suite.userA, "")
	tp2, _ := suite.ts.GetNewTokenPair(suite.userA, "")

	suite.Assert().NoError(suite.ts.DeleteAllUserTokens(suite.userA))
	_, err := suite.ts.GetNewTokenPair(suite.userA, tp1.RefreshToken.ID.String())
	suite.Require().Error(err)

	_, err = suite.ts.GetNewTokenPair(suite.userA, tp2.RefreshToken.ID.String())
	suite.Require().Error(err)

	suite.Assert().NoError(suite.ts.DeleteAllUserTokens(suite.userA))
}
