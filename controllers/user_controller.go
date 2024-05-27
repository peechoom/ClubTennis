package controllers

import (
	"ClubTennis/models"
	"ClubTennis/services"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userservice  *services.UserService
	matchservice *services.MatchService
}

func NewUserController(userService *services.UserService, matchService *services.MatchService) *UserController {
	return &UserController{
		userservice:  userService,
		matchservice: matchService,
	}
}

//----------------------------------------------------------------------------------------------------------------
// GET handlers
/*
	GET ../members

gets all members. a quite heavy operation
*/
func (ctrl *UserController) GetAllMembers(c *gin.Context) {
	users, err := ctrl.userservice.FindAll()
	if err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, "record not found")
		return
	}
	c.JSON(http.StatusOK, users)
}

/*
	GET ../members/{id}

gets the member with the given ID. Can either be their numeric ID in the db or their unity ID
*/
func (ctrl *UserController) GetMemberByID(c *gin.Context) {
	fetch := c.Param("id")
	if len(fetch) == 0 {
		c.String(http.StatusBadRequest, "id not provided")
		return
	}

	ui, id := getIDorUnityID(fetch)
	var user *models.User
	var err error

	if ui != 0 { //try to use numeric id
		user, err = ctrl.userservice.FindByID(ui)
		if err != nil || user == nil {
			c.String(http.StatusNotFound, "user not found")
			return
		}
	} else { //try to use unity id
		if !validateUnityID(id) {
			c.String(http.StatusBadRequest, "not a valid Unity/SKEMA ID")
			return
		}
		user, err = ctrl.userservice.FindByUnityID(id)
		if err != nil || user == nil {
			c.String(http.StatusNotFound, "user not found")
			return
		}
	}
	c.JSON(http.StatusOK, user)
}

//---------------------------------------------------------------------------------------------------------------
// POST handlers
/*
	POST ../members
	expects json containing fields: UnityID, FirstName, LastName, Affiliation, Email

creates a new member. responds with json representing the member. only admins should be able to call this
*/
func (ctrl *UserController) CreateMember(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil || len(user.UnityID) == 0 {
		c.String(http.StatusBadRequest, "error interpreting passed user")
		return
	}
	user.ID = 0

	err = ctrl.userservice.Save(&user)
	if err != nil || user.ID == 0 {
		c.String(http.StatusBadRequest, "error saving user")
		return
	}
	c.JSON(http.StatusCreated, user)
}

//---------------------------------------------------------------------------------------------------------------
// PUT handlers
/*
	PUT ../members/{id}
	expects JSON representing the modified user.

edits the member with the given ID. For changing things like name, email etc etc etc
*/
func (ctrl *UserController) EditMember(c *gin.Context) {
	var nu models.User
	var err error
	parse, err := strconv.ParseUint(c.Param("id"), 10, 0)

	if err != nil || parse == 0 {
		c.String(http.StatusBadRequest, "illegal id provided")
		return
	}
	err = c.BindJSON(&nu)
	if err != nil {
		c.String(http.StatusBadRequest, "error interpreting passed user")
		return
	}

	nu.ID = uint(parse)

	fetched, err := ctrl.userservice.FindByID(nu.ID)
	if err != nil {
		c.String(http.StatusNotFound, "user not found")
		return
	}
	fetched.EditUser(&nu)
	err = ctrl.userservice.Save(fetched)
	if err != nil {
		c.String(http.StatusInternalServerError, "error saving user")
		return
	}
	c.String(http.StatusOK, "user edited successfully")
}

//---------------------------------------------------------------------------------------------------------------
// DELETE handlers
/*
	DELETE ../members/{id}

deletes the user with the given ID from the database. only admins should be able to call this. can be unity or numeric
*/
func (ctrl *UserController) DeleteMember(c *gin.Context) {
	fetch := c.Param("id")
	if len(fetch) == 0 {
		c.String(http.StatusBadRequest, "id not provided")
		return
	}

	ui, unity := getIDorUnityID(fetch)
	var err error

	if ui != 0 { //numeric id
		err = ctrl.userservice.DeleteByID(ui)
		if err != nil {
			c.String(http.StatusNotFound, "user not deleted: record not found")
			return
		}
	} else { //unity id
		if !validateUnityID(unity) {
			c.String(http.StatusBadRequest, "poorly formatted unityID")
		}
		err = ctrl.userservice.DeleteByUnityID(unity)
		if err != nil {
			c.String(http.StatusNotFound, "user not deleted: record not found")
			return
		}
	}
	c.String(http.StatusOK, "user deleted successfuly")
}

//---------------------------------------------------------------------------------------------------------------
// UTILITY functions

// returns the id either parsed as a uint or a string. the other value will be 0 or ""
func getIDorUnityID(id string) (uint, string) {
	num, err := strconv.ParseUint(id, 10, 0)
	if err != nil || num == 0 {
		return 0, id
	}
	return uint(num), ""
}

// returns true if the unityID is valid. should also work for skema
func validateUnityID(input string) bool {
	regexPattern := "^[a-zA-Z0-9.]+$"
	matched, _ := regexp.MatchString(regexPattern, input)
	return matched
}
