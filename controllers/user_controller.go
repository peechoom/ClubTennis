package controllers

import (
	"ClubTennis/models"
	"ClubTennis/services"
	"net/http"
	"regexp"
	"strconv"
	"sync"

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

gets all members, only matches in the last 14 days are preloaded and all players are marked
challengeable or not relative to the signed in principal
*/
func (ctrl *UserController) GetAllMembers(c *gin.Context) {
	users, err := ctrl.userservice.FindAll()

	if err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, "record not found")
		return
	}

	uid := c.GetUint("user_id")
	if uid == 0 {
		//no user signed in, but authorized
		c.JSON(http.StatusOK, users)
		return
	}
	var p *models.User
	for i := range users {
		if users[i].ID == uid {
			p = &users[i]
			break
		}
	}
	if p == nil {
		//uid in list but not found... lets just return
		c.JSON(http.StatusInternalServerError, gin.H{"error": "userID provided but not found"})
		return
	}

	var wg sync.WaitGroup
	const WORKER_COUNT int = 6 // each user should get idk 6 threads. That makes what like 10 users per thread
	userChannel := make(chan *models.User)

	worker := func() {
		defer wg.Done()
		for u := range userChannel {
			u.IsChallengeable, _ = p.CanChallenge(u)
		}
	}
	for i := 0; i < WORKER_COUNT; i++ {
		wg.Add(1)
		go worker()
	}
	go func() {
		for i := range users {
			userChannel <- &users[i]
		}
		close(userChannel)
	}()

	wg.Wait()

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

/*
	GET ../cutoff/:ladder

returns the cutoff for the mens or womens red team. either m or w
*/
func (ctrl *UserController) GetCutoff(c *gin.Context) {
	ladder := c.Param("ladder")
	var num int
	if ladder == "m" {
		num = ctrl.userservice.GetLadderCutoff(models.MENS_LADDER)
	} else if ladder == "w" {
		num = ctrl.userservice.GetLadderCutoff(models.WOMENS_LADDER)
	} else {
		c.String(http.StatusBadRequest, "not a valid ladder")
		return
	}
	if num == -1 {
		c.String(http.StatusInternalServerError, "error fetching num")
		return
	}
	c.String(http.StatusOK, strconv.FormatInt(int64(num), 10))
}

//---------------------------------------------------------------------------------------------------------------
// POST handlers
/*
	POST ../members
	expects json containing fields: UnityID, FirstName, LastName, Affiliation, SigninEmail, ContactEmail, ladder

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

/*
POST ../cutoff/:ladder

sets the cutoff ofthe ladder, either m or w. expects JSON field with "cutoff": num
*/
func (ctrl *UserController) SetCutoff(c *gin.Context) {
	ladder := c.Param("ladder")
	var payload map[string]interface{}
	c.BindJSON(&payload)
	val, ok := payload["cutoff"]
	if !ok {
		c.String(http.StatusBadRequest, "cutoff field not present")
		return
	}
	num := int(val.(float64))
	if ladder == "m" {
		if err := ctrl.userservice.SetLadderCutoff(models.MENS_LADDER, num); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	} else if ladder == "w" {
		if err := ctrl.userservice.SetLadderCutoff(models.MENS_LADDER, num); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		c.String(http.StatusBadRequest, "not a valid ladder")
		return
	}
	c.String(http.StatusOK, "cutoff set")
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
