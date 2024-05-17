package repositories_test

import (
	"ClubTennis/config"
	"ClubTennis/initializers"
	"ClubTennis/repositories"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TokenRepoTestSuite struct {
	suite.Suite
	repo *repositories.TokenRepository
}

func (suite *TokenRepoTestSuite) SetupTest() {
	c, err := config.LoadConfig("/home/alec/go/src/ClubTennis/config/config.json")
	if err != nil || c == nil {
		panic(err.Error())
	}
	client := initializers.GetTestClient(c)
	if client == nil {
		panic("couldn't get client")
	}
	suite.repo = repositories.NewTokenRepository(client)
	if suite.repo == nil {
		panic("couldnt get client")
	}
}

func TestTokenRepoSuite(t *testing.T) {
	suite.Run(t, new(TokenRepoTestSuite))
}

func (suite *TokenRepoTestSuite) TestSetDeleteToken() {
	var userID string = "44"
	var tokenID uuid.UUID
	tokenID, _ = uuid.NewRandom()

	err := suite.repo.SetRefreshToken(userID, tokenID.String(), time.Hour)
	suite.Require().NoError(err)

	err = suite.repo.DeleteRefreshToken(userID, tokenID.String())
	suite.Require().NoError(err)

	err = suite.repo.DeleteRefreshToken(userID, tokenID.String())
	suite.Require().Error(err)
}

func (suite *TokenRepoTestSuite) TestDeleteAll() {
	var userID string = "44"
	tokenID1, _ := uuid.NewRandom()
	err := suite.repo.SetRefreshToken(userID, tokenID1.String(), time.Hour)
	suite.Require().NoError(err)

	tokenID2, _ := uuid.NewRandom()
	err = suite.repo.SetRefreshToken(userID, tokenID2.String(), time.Hour)
	suite.Require().NoError(err)

	tokenID3, _ := uuid.NewRandom()
	err = suite.repo.SetRefreshToken(userID, tokenID3.String(), time.Hour)
	suite.Require().NoError(err)

	tokenID4, _ := uuid.NewRandom()
	err = suite.repo.SetRefreshToken(userID, tokenID4.String(), time.Hour)
	suite.Require().NoError(err)

	err = suite.repo.DeleteUserRefreshTokens(userID)
	suite.Require().NoError(err)

	err = suite.repo.DeleteRefreshToken(userID, tokenID1.String())
	suite.Require().Error(err)
	err = suite.repo.DeleteRefreshToken(userID, tokenID2.String())
	suite.Require().Error(err)
	err = suite.repo.DeleteRefreshToken(userID, tokenID3.String())
	suite.Require().Error(err)
	err = suite.repo.DeleteRefreshToken(userID, tokenID4.String())
	suite.Require().Error(err)
}
