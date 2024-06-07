package controllers

import (
	"ClubTennis/models"
	"ClubTennis/services"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const ACTIVE_QUERY string = "active"
const PLAYERS_QUERY string = "player"

// how "recent" a "recent" match is
const RECENT_TIME_SPAN = 7 * time.Hour * 24

type MatchController struct {
	matchservice *services.MatchService
	userservice  *services.UserService
	emailservice *services.EmailService
}

func NewMatchController(matchService *services.MatchService, userService *services.UserService, emailService *services.EmailService) *MatchController {
	return &MatchController{
		matchservice: matchService,
		userservice:  userService,
		emailservice: emailService,
	}
}

//---------------------------------------------------------------------------------------------------------
// GET HANDLERS

/*
	GET .../matches

returns all matches in the db
*/
func (ctrl *MatchController) GetAllMatches(c *gin.Context) {
	matches, err := ctrl.matchservice.FindAll()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, matches)
}

/*
	GET .../matches/{id}

returns the match with the given ID
*/
func (ctrl *MatchController) GetMatchByID(c *gin.Context) {
	id, err := getID(c)
	if err != nil || id == 0 {
		c.String(http.StatusBadRequest, "ID is nan")
		return
	}

	fetched, err := ctrl.matchservice.FindByID(id)

	if err != nil || fetched == nil {
		c.String(http.StatusNotFound, "record not found")
		return
	}

	c.JSON(http.StatusOK, fetched)

}

/*
	GET .../matches?active=true&player=1234&player=345

returns matches with a query string containing 1 or 0 active queries and 0 or many player queries. THis endpoint is hideous
*/
func (ctrl *MatchController) GetMatch(c *gin.Context) {
	var query url.Values
	var err error

	query, err = url.ParseQuery(c.Request.URL.RawQuery)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	//if no params in the query string, return everything
	if len(query) == 0 {
		ctrl.GetAllMatches(c)
		return
	}
	if !query.Has(PLAYERS_QUERY) && !query.Has(ACTIVE_QUERY) {
		c.Status(http.StatusBadRequest)
		return
	}
	active, err := isActive(query)
	//active query dne
	if err != nil {
		ls, err := strsToUints(query[PLAYERS_QUERY])
		if err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}
		fetched, err := ctrl.matchservice.FindByPlayerID(ls...)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, fetched)
		return
	}

	if query.Has(PLAYERS_QUERY) {
		ls, err := strsToUints(query[PLAYERS_QUERY])
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		fetched, err := ctrl.matchservice.FindByPlayerIDAndActive(active, ls...)
		if err != nil || fetched == nil || len(fetched) == 0 {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, fetched)
		return
	}
	fetched, err := ctrl.matchservice.FindByActive(active)
	if err != nil || fetched == nil || len(fetched) == 0 {
		c.String(http.StatusNotFound, "value not found")
		return
	}
	c.JSON(http.StatusOK, fetched)
}

/*
	GET .../matches/recent

returns all recent matches. player ids included but not player data. responsibility of frontend
*/
func (ctrl *MatchController) GetRecentMatches(c *gin.Context) {
	matches, err := ctrl.matchservice.FindAllRecentMatches(RECENT_TIME_SPAN)
	if err != nil {
		c.Error(err)
		log.Print(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "could not find recent matches"})
		return
	}
	c.JSON(http.StatusOK, matches)
}

//----------------------------------------------------------------------------------------------------------------
// POST HANDLERS

/*
POST .../matches
expects: challengerID and challengedID in json
*/
func (ctrl *MatchController) Challenge(c *gin.Context) {

	var m map[string]int
	err := c.ShouldBindBodyWithJSON(&m)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "post form badly formatted")
		return
	}
	challengerID := m["challengerID"]
	challengedID := m["challengedID"]

	if challengerID == 0 || challengedID == 0 {
		log.Println(challengerID)
		log.Println(challengedID)
		c.String(http.StatusBadRequest, "post form badly formatted")
		return
	}

	principleID := c.GetUint("user_id")
	if principleID != uint(challengerID) {
		c.String(http.StatusBadRequest, "You are not the challenger")
	}

	challenger, err := ctrl.userservice.FindByID(uint(challengerID))
	if err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, "Challenger ID not found")
		return
	}
	challenged, err := ctrl.userservice.FindByID(uint(challengedID))
	if err != nil {
		c.Error(err)
		c.String(http.StatusNotFound, "Opponent ID not found")
		return
	}

	match, err := challenger.Challenge(challenged)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = ctrl.matchservice.Save(match)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Internal error")
		return
	}
	e := ctrl.notifyPlayers(c, match)

	if e != nil {
		if e != errChallengerNotNotified {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Challenge created but emails failed to generate. Please contact %s to let %s know about the challenge.", challenged.Email, challenged.FirstName))
		} else {
			c.String(http.StatusAccepted, fmt.Sprintf("Challenge created successfully, but due to an internal server error you will not be emailed a reciept. %s has been sucessfully contacted, and should contact you soon", challenged.FirstName))
		}
	} else {
		c.JSON(http.StatusCreated, *match)
	}
}

var errChallengedNotNotified error = errors.New("challenged player not notified")
var errChallengerNotNotified error = errors.New("challenger not notified")

func (ctrl *MatchController) notifyPlayers(c *gin.Context, match *models.Match) error {
	challengerEmail, challengedEmail := ctrl.emailservice.MakeChallengeEmails(match.Challenger(), match.Challenged())
	if challengerEmail == nil || challengedEmail == nil {
		e := errors.New("could not create emails")
		c.Error(e)
		return e
	}
	if ctrl.emailservice.Send(challengedEmail) != nil {
		c.Error(errChallengedNotNotified)
		return errChallengedNotNotified
	}
	if ctrl.emailservice.Send(challengerEmail) != nil {
		c.Error(errChallengerNotNotified)
		return errChallengerNotNotified
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------
// PATCH HANDLERS

/*
	PATCH .../matches/{id}
	expects: challengerScore and challengedScore in json

endpoint for submitting the score of a match once it has concluded
*/
func (ctrl *MatchController) SubmitScore(c *gin.Context) {
	id, err := getID(c)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	match, err := ctrl.matchservice.FindByID(id)
	if err != nil {
		c.Error(err)
		return
	}
	if match == nil {
		c.Status(http.StatusBadRequest)
		return
	}

	var m map[string]int
	err = c.ShouldBindBodyWithJSON(&m)
	if err != nil {
		log.Println(err)
		c.String(http.StatusBadRequest, "post form badly formatted")
		return
	}
	challengerScore := m["challengerScore"]
	challengedScore := m["challengedScore"]

	err = match.SubmitScore(int(challengerScore), int(challengedScore))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = ctrl.matchservice.Save(match)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, match)
}

// ----------------------------------------------------------------------------------------------------------------
// UTILITY functions

func isActive(query url.Values) (bool, error) {
	aq := query[ACTIVE_QUERY]
	if len(aq) != 1 {
		return false, errors.New("multiple active values found")
	}
	isActive, err := strconv.ParseBool(aq[0])
	if err != nil {
		return false, errors.New("bad formatting on active value")
	}
	return isActive, nil
}

func strsToUints(strings []string) (uints []uint, err error) {
	for _, str := range strings {
		u, err := strconv.ParseUint(str, 10, 0)
		if err != nil {
			return nil, err
		}
		uints = append(uints, (uint(u)))
	}
	return uints, nil
}
func getID(c *gin.Context) (uint, error) {
	idString := c.Param("id")
	parse, err := strconv.ParseUint(idString, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(parse), nil
}
