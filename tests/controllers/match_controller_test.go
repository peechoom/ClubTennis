package controllers_test

import (
	"ClubTennis/controllers"
	"ClubTennis/initializers"
	"ClubTennis/models"
	"ClubTennis/services"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type MatchControllerTestSuite struct {
	suite.Suite
	router *gin.Engine
	w      *httptest.ResponseRecorder
	ctrl   *controllers.MatchController
	ms     *services.MatchService
	us     *services.UserService
	userA  *models.User
	userB  *models.User
	userC  *models.User
	userD  *models.User
	userE  *models.User
}

func (suite *MatchControllerTestSuite) SetupTest() {
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
	s := services.SetupServices(db)
	suite.ctrl = controllers.NewMatchController(s.MatchService, s.UserService)
	suite.ms = s.MatchService
	suite.us = s.UserService

	suite.router = initializers.GetTestEngine()
	suite.w = httptest.NewRecorder()

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

func TestMatchControllerTestSuite(t *testing.T) {
	suite.Run(t, new(MatchControllerTestSuite))
}

func (suite *MatchControllerTestSuite) TestNewChallenge() {
	userA, _ := suite.us.FindByUnityID(suite.userA.UnityID)
	userC, _ := suite.us.FindByUnityID(suite.userC.UnityID)

	req, _ := http.NewRequest("POST", "/club/matches", nil)
	req.Header.Set("Content-Type", "application/json")
	var v url.Values = map[string][]string{
		"challengerID": {strconv.Itoa(int(userC.ID))},
		"challengedID": {strconv.Itoa(int(userA.ID))},
	}
	req.PostForm = v

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusCreated, suite.w.Code)

	matches, _ := suite.ms.FindAll()
	suite.Require().Len(matches, 1)

	suite.Require().Equal(userC.ID, matches[0].ChallengerID)
	suite.Require().Equal(userA.ID, matches[0].ChallengedID)
	suite.Require().True(matches[0].IsActive)

}

func (suite *MatchControllerTestSuite) TestNewChallengeIllegal() {
	userA, _ := suite.us.FindByUnityID(suite.userA.UnityID)
	userD, _ := suite.us.FindByUnityID(suite.userD.UnityID)

	req, _ := http.NewRequest("POST", "/club/matches", nil)
	req.Header.Set("Content-Type", "application/json")
	var v url.Values = map[string][]string{
		"challengerID": {strconv.Itoa(int(userD.ID))},
		"challengedID": {strconv.Itoa(int(userA.ID))},
	}
	req.PostForm = v

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusForbidden, suite.w.Code)

	matches, _ := suite.ms.FindAll()
	suite.Require().Len(matches, 0)

	req, _ = http.NewRequest("POST", "/club/matches", nil)
	req.Header.Set("Content-Type", "application/json")
	v = map[string][]string{
		"challengerID": {strconv.Itoa(int(userA.ID))},
		"challengedID": {strconv.Itoa(int(userA.ID))},
	}
	req.PostForm = v

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusForbidden, suite.w.Code)

	matches, _ = suite.ms.FindAll()
	suite.Require().Len(matches, 0)
}

func (suite *MatchControllerTestSuite) TestSubmitMatchScore() {
	userA, _ := suite.us.FindByUnityID(suite.userA.UnityID)
	userC, _ := suite.us.FindByUnityID(suite.userC.UnityID)

	req, _ := http.NewRequest("POST", "/club/matches", nil)
	req.Header.Set("Content-Type", "application/json")
	var v url.Values = map[string][]string{
		"challengerID": {strconv.Itoa(int(userC.ID))},
		"challengedID": {strconv.Itoa(int(userA.ID))},
	}
	req.PostForm = v

	suite.router.ServeHTTP(suite.w, req)
	suite.Assert().Equal(http.StatusCreated, suite.w.Code)

	//submit the score
	userA, _ = suite.us.FindByUnityID(suite.userA.UnityID)
	userC, _ = suite.us.FindByUnityID(suite.userC.UnityID)

	match, err := suite.ms.FindByPlayerIDAndActive(true, userA.ID, userC.ID)
	suite.Assert().NoError(err)
	suite.Assert().Len(match, 1)

	route := fmt.Sprintf("/club/matches/%d", match[0].ID)
	req, _ = http.NewRequest("PATCH", route, nil)
	req.Header.Set("Content-Type", "application/json")
	v = map[string][]string{
		"challengerScore": {"6"},
		"challengedScore": {"2"},
	}
	req.PostForm = v
	suite.w = httptest.NewRecorder()

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusOK, suite.w.Code)

	matches, _ := suite.ms.FindAll()
	suite.Assert().Len(matches, 1)

	suite.Require().False(matches[0].IsActive)
	suite.Require().Equal(models.EncodeScore(6, 2), matches[0].Score)
	winnerID, _ := matches[0].Winner()
	suite.Require().Equal(userC.ID, winnerID)

}

func (suite *MatchControllerTestSuite) TestSubmitMatchScoreIllegal() {
	userA, _ := suite.us.FindByUnityID(suite.userA.UnityID)
	userC, _ := suite.us.FindByUnityID(suite.userC.UnityID)

	req, _ := http.NewRequest("POST", "/club/matches", nil)
	req.Header.Set("Content-Type", "application/json")
	var v url.Values = map[string][]string{
		"challengerID": {strconv.Itoa(int(userC.ID))},
		"challengedID": {strconv.Itoa(int(userA.ID))},
	}
	req.PostForm = v

	suite.router.ServeHTTP(suite.w, req)
	suite.Assert().Equal(http.StatusCreated, suite.w.Code)

	//submit the score
	userA, _ = suite.us.FindByUnityID(suite.userA.UnityID)
	userC, _ = suite.us.FindByUnityID(suite.userC.UnityID)

	match, err := suite.ms.FindByPlayerIDAndActive(true, userA.ID, userC.ID)
	suite.Assert().NoError(err)
	suite.Assert().Len(match, 1)

	route := fmt.Sprintf("/club/matches/%d", match[0].ID)
	req, _ = http.NewRequest("PATCH", route, nil)
	req.Header.Set("Content-Type", "application/json")
	v = map[string][]string{
		"challengerScore": {"6"},
		"challengedScore": {"6"},
	}
	req.PostForm = v
	suite.w = httptest.NewRecorder()

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusBadRequest, suite.w.Code)

	matches, _ := suite.ms.FindAll()
	fetched := matches[0]
	suite.Require().True(fetched.IsActive)
	suite.Require().Zero(fetched.Score)

}

func (suite *MatchControllerTestSuite) TestGetMatches() {
	//save a match to the db
	match, _ := suite.userC.Challenge(suite.userA)
	suite.ms.Save(match)

	route := fmt.Sprintf("/club/matches/%d", match.ID)
	req, _ := http.NewRequest("GET", route, nil)

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusOK, suite.w.Code)

	suite.Require().Contains(suite.w.Body.String(), "\"ChallengerID\":"+strconv.Itoa(int(suite.userC.ID)))
	suite.Require().Contains(suite.w.Body.String(), "\"ChallengedID\":"+strconv.Itoa(int(suite.userA.ID)))

	suite.w = httptest.NewRecorder()
	route = fmt.Sprintf("/club/matches/%d", match.ID+5)
	req, _ = http.NewRequest("GET", route, nil)

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusNotFound, suite.w.Code)

	suite.w = httptest.NewRecorder()
	route = "/club/matches/asdf"
	req, _ = http.NewRequest("GET", route, nil)

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusBadRequest, suite.w.Code)

}

func (suite *MatchControllerTestSuite) TestGetMatch() {
	//save a match1 to the db
	match1, _ := suite.userC.Challenge(suite.userA)
	match2, _ := suite.userD.Challenge(suite.userB)
	suite.ms.Save(match1)
	suite.ms.Save(match2)

	route := fmt.Sprintf("/club/matches?active=true&player=%d&player=%d", match1.ChallengerID, match1.ChallengedID)
	req, _ := http.NewRequest("GET", route, nil)

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusOK, suite.w.Code)

	suite.Require().Contains(suite.w.Body.String(), "\"ChallengerID\":"+strconv.Itoa(int(suite.userC.ID)))
	suite.Require().Contains(suite.w.Body.String(), "\"ChallengedID\":"+strconv.Itoa(int(suite.userA.ID)))
	suite.Require().NotContains(suite.w.Body.String(), "\"ChallengerID\":"+strconv.Itoa(int(suite.userD.ID)))
	suite.Require().NotContains(suite.w.Body.String(), "\"ChallengedID\":"+strconv.Itoa(int(suite.userB.ID)))

	route = fmt.Sprintf("/club/matches?active=true&player=%d", match1.ChallengerID)
	req, _ = http.NewRequest("GET", route, nil)
	suite.w = httptest.NewRecorder()

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusOK, suite.w.Code)

	suite.Require().Contains(suite.w.Body.String(), "\"ChallengerID\":"+strconv.Itoa(int(suite.userC.ID)))
	suite.Require().Contains(suite.w.Body.String(), "\"ChallengedID\":"+strconv.Itoa(int(suite.userA.ID)))
	suite.Require().NotContains(suite.w.Body.String(), "\"ChallengerID\":"+strconv.Itoa(int(suite.userD.ID)))
	suite.Require().NotContains(suite.w.Body.String(), "\"ChallengedID\":"+strconv.Itoa(int(suite.userB.ID)))

	route = "/club/matches?active=true"
	req, _ = http.NewRequest("GET", route, nil)
	suite.w = httptest.NewRecorder()

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusOK, suite.w.Code)

	suite.Require().Contains(suite.w.Body.String(), "\"ChallengerID\":"+strconv.Itoa(int(suite.userC.ID)))
	suite.Require().Contains(suite.w.Body.String(), "\"ChallengedID\":"+strconv.Itoa(int(suite.userA.ID)))
	suite.Require().Contains(suite.w.Body.String(), "\"ChallengerID\":"+strconv.Itoa(int(suite.userD.ID)))
	suite.Require().Contains(suite.w.Body.String(), "\"ChallengedID\":"+strconv.Itoa(int(suite.userB.ID)))

	route = "/club/matches?active=false"
	req, _ = http.NewRequest("GET", route, nil)
	suite.w = httptest.NewRecorder()

	suite.router.ServeHTTP(suite.w, req)
	suite.Assert().Equal(http.StatusNotFound, suite.w.Code)

	suite.Require().NotContains(suite.w.Body.String(), "\"ChallengerID\":"+strconv.Itoa(int(suite.userC.ID)))
	suite.Require().NotContains(suite.w.Body.String(), "\"ChallengedID\":"+strconv.Itoa(int(suite.userA.ID)))
	suite.Require().NotContains(suite.w.Body.String(), "\"ChallengerID\":"+strconv.Itoa(int(suite.userD.ID)))
	suite.Require().NotContains(suite.w.Body.String(), "\"ChallengedID\":"+strconv.Itoa(int(suite.userB.ID)))

}
