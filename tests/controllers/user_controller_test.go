package controllers

import (
	"ClubTennis/controllers"
	"ClubTennis/initializers"
	"ClubTennis/models"
	"ClubTennis/services"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type UserControllerTestSuite struct {
	suite.Suite
	router *gin.Engine
	w      *httptest.ResponseRecorder
	ctrl   *controllers.UserController
	ms     *services.MatchService
	us     *services.UserService
	userA  *models.User
	userB  *models.User
	userC  *models.User
	userD  *models.User
	userE  *models.User
}

func (suite *UserControllerTestSuite) SetupTest() {
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
	suite.ctrl = controllers.NewUserController(s.UserService, s.MatchService)
	suite.ms = s.MatchService
	suite.us = s.UserService

	suite.router = initializers.GetTestEngine()
	suite.w = httptest.NewRecorder()

	suite.userA, _ = models.NewUser("shboil4", "ncsu", "Sam", "Boiland", "shboil4@ncsu.edu", models.MENS_LADDER)
	suite.userB, _ = models.NewUser("jbeno5", "ncsu", "James", "Benolli", "jbeno5@ncsu.edu", models.MENS_LADDER)
	suite.userC, _ = models.NewUser("pdiddy4", "ncsu", "Puff", "Daddy", "pdiddy@ncsu.edu", models.MENS_LADDER)
	suite.userD, _ = models.NewUser("jobitch2", "ncsu", "Joel", "Embitch", "jobitch@ncsu.edu", models.MENS_LADDER)
	suite.userE, _ = models.NewUser("myprince2", "ncsu", "Lebron", "James", "myprince2@ncsu.edu", models.MENS_LADDER)

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

func TestUserControllerSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

func (suite *UserControllerTestSuite) TestGetUserByID() {
	//correct unity id
	var route string = "/club/members/pdiddy4"
	req, _ := http.NewRequest("GET", route, nil)
	c, _ := gin.CreateTestContext(suite.w) //idk how to make ts work
	c.Set("user_id", uint(1))

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusOK, suite.w.Code)

	var fetched models.User
	json.Unmarshal(suite.w.Body.Bytes(), &fetched)

	suite.Require().Equal(*suite.userC, fetched)

	//DNE unity id
	route = "/club/members/pdiddy3"
	req, _ = http.NewRequest("GET", route, nil)
	suite.w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(suite.w)
	c.Set("user_id", uint(1))
	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusNotFound, suite.w.Code)

	//numerical id
	id := suite.userC.ID
	route = fmt.Sprintf("/club/members/%d", id)
	req, _ = http.NewRequest("GET", route, nil)
	suite.w = httptest.NewRecorder()

	suite.router.ServeHTTP(suite.w, req)
	fetched = models.User{}
	json.Unmarshal(suite.w.Body.Bytes(), &fetched)

	suite.Require().Equal(*suite.userC, fetched)

	//DNE numerical id
	id = 10000
	route = fmt.Sprintf("/club/members/%d", id)
	req, _ = http.NewRequest("GET", route, nil)
	suite.w = httptest.NewRecorder()

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusNotFound, suite.w.Code)

}

func (suite *UserControllerTestSuite) TestCreateNewUser() {
	var route string = "/admin/members"
	newGuy, _ := models.NewUser("kwest4", "ncsu", "Kanye", "West", "kwest4@ncsu.edu", models.MENS_LADDER)
	marsh, _ := json.Marshal(newGuy)
	req, _ := http.NewRequest("POST", route, bytes.NewBuffer(marsh))
	req.Header.Set("Content-Type", "application/json")

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusCreated, suite.w.Code)
	var fetched models.User
	json.Unmarshal(suite.w.Body.Bytes(), &fetched)

	suite.Require().Equal(newGuy.UnityID, fetched.UnityID)
	suite.Require().NotZero(fetched.ID)
	suite.Require().NotZero(fetched.CreatedAt.Unix())
}

func (suite *UserControllerTestSuite) TestEditUser() {
	userE, _ := suite.us.FindByID(suite.userE.ID)
	route := fmt.Sprintf("/admin/members/%d", userE.ID)
	userE.UnityID = "wewoo1"
	userE.LastName = "wow"
	userE.FirstName = ""
	userE.Email = "newEmail@email.com"

	marsh, _ := json.Marshal(userE)

	req, _ := http.NewRequest("PUT", route, bytes.NewBuffer(marsh))

	suite.router.ServeHTTP(suite.w, req)

	suite.Require().Equal(http.StatusOK, suite.w.Code)

	fetched1, _ := suite.us.FindByID(userE.ID)
	fetched2, _ := suite.us.FindByUnityID("wewoo1")

	suite.Require().Equal(fetched1, fetched2)
	suite.Require().Equal("wow", fetched1.LastName)
	suite.Require().Equal("Lebron", fetched1.FirstName)
	suite.Require().Equal("newEmail@email.com", fetched1.Email)

	suite.w = httptest.NewRecorder()
	route = "/admin/members/123412341234"
	req, _ = http.NewRequest("PUT", route, bytes.NewBuffer(marsh))
	suite.router.ServeHTTP(suite.w, req)

	suite.Require().Equal(http.StatusNotFound, suite.w.Code)

	suite.w = httptest.NewRecorder()
	route = "/admin/members/asdfasdfadsf"
	req, _ = http.NewRequest("PUT", route, bytes.NewBuffer(marsh))
	suite.router.ServeHTTP(suite.w, req)

	suite.Require().Equal(http.StatusBadRequest, suite.w.Code)
}

func (suite *UserControllerTestSuite) TestDeleteUser() {
	route := fmt.Sprintf("/admin/members/%d", suite.userE.ID)
	req, _ := http.NewRequest("DELETE", route, nil)

	suite.router.ServeHTTP(suite.w, req)
	suite.Require().Equal(http.StatusOK, suite.w.Code)

	fetched, err := suite.us.FindByID(suite.userE.ID)
	suite.Assert().Error(err)
	suite.Require().Nil(fetched)

}
