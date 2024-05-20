package repositories_test

import (
	"ClubTennis/config"
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
	err := config.LoadConfig("/home/alec/go/src/ClubTennis/config/.env")
	if err != nil {
		panic(err.Error())
	}
	suite.repo = repositories.NewTokenRepository()
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
	var err error
	tokenID, _ = uuid.NewRandom()

	suite.repo.SetRefreshToken(userID, tokenID.String(), time.Hour)
	suite.Require().NoError(err)

	err = suite.repo.DeleteRefreshToken(userID, tokenID.String())
	suite.Require().NoError(err)

	err = suite.repo.DeleteRefreshToken(userID, tokenID.String())
	suite.Require().Error(err)
}
