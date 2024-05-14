package controllers

import (
	"ClubTennis/services"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const ACTIVE_QUERY string = "active"
const PLAYERS_QUERY string = "player"

type MatchController struct {
	matchservice *services.MatchService
	userservice  *services.UserService
}

func NewMatchController(db *gorm.DB) *MatchController {
	return &MatchController{
		matchservice: services.NewMatchService(db),
		userservice:  services.NewUserService(db),
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
		if err != nil || fetched == nil || len(fetched) == 0 {
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

//----------------------------------------------------------------------------------------------------------------
// POST HANDLERS

/*
POST .../matches
expects: challengerID and challengedID in json
*/
func (ctrl *MatchController) Challenge(c *gin.Context) {

	challengerIDstr := c.PostForm("challengerID")
	challengedIDstr := c.PostForm("challengedID")
	if len(challengerIDstr) == 0 || len(challengedIDstr) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	challengerID, err1 := strconv.ParseUint(challengerIDstr, 10, 0)
	challengedID, err2 := strconv.ParseUint(challengedIDstr, 10, 0)

	if err1 != nil || err2 != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	challenger, err := ctrl.userservice.FindByID(uint(challengerID))
	if err != nil {
		c.Error(err)
		return
	}
	challenged, err := ctrl.userservice.FindByID(uint(challengedID))
	if err != nil {
		c.Error(err)
		return
	}

	match, err := challenger.Challenge(challenged)
	if err != nil {
		c.String(http.StatusForbidden, err.Error())
		return
	}

	err = ctrl.matchservice.Save(match)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, *match)
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

	challengerScoreStr := c.PostForm("challengerScore")
	challengedScoreStr := c.PostForm("challengedScore")
	if len(challengedScoreStr) == 0 || len(challengedScoreStr) == 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	challengerScore, err1 := strconv.ParseInt(challengerScoreStr, 10, 0)
	challengedScore, err2 := strconv.ParseInt(challengedScoreStr, 10, 0)
	if err1 != nil || err2 != nil {
		c.Status(http.StatusBadRequest)
		return
	}

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
