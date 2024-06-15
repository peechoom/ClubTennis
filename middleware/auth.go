package middleware

import (
	"ClubTennis/models"
	"ClubTennis/services"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Authenticator struct {
	tokenService *services.TokenService
	userService  *services.UserService
	host         string
	cacheMutex   *sync.RWMutex
	writing      bool
	adminIDs     map[uint]bool
}

func NewAuthenticator(tokenService *services.TokenService, userService *services.UserService, host string) *Authenticator {
	a := &Authenticator{
		tokenService: tokenService,
		userService:  userService,
		host:         host,
		cacheMutex:   &sync.RWMutex{},
		writing:      false,
	}
	a.resetAdminCache()

	return a
}

func DoNotCache(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
}

func (a *Authenticator) AuthenticateMember(c *gin.Context) {
	var err error
	idTokenString, err := c.Cookie("id_token")
	if err != nil && err != http.ErrNoCookie {
		c.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var userID uint
	userID, err = a.tokenService.ValidateIDToken(idTokenString)
	if err == nil {
		//id token valid, user is authenticated
		//set the userID in this context so the principal can be used later
		c.Set("user_id", userID)
		c.Next()
		return
	}
	userID, err = a.cycleRefreshTokens(c)
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("user_id", userID)
	c.Next()
}

func (a *Authenticator) AuthenticateAdmin(c *gin.Context) {
	var err error
	idTokenString, err := c.Cookie("id_token")
	if err != nil && err != http.ErrNoCookie {
		c.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var userID uint
	userID, err = a.tokenService.ValidateIDToken(idTokenString)
	if err == nil {

		a.cacheMutex.RLock()
		cond := a.adminIDs[userID]
		a.cacheMutex.RUnlock()

		//verify admin status. root user has uid 0
		if userID != 0 && !cond {
			//check database, update cache if an update has been made.
			u, err := a.userService.FindByID(userID)
			if err != nil || !u.IsOfficer {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}
		go a.resetAdminCache()

		c.Set("user_id", userID)
		c.Next()
		return
	}

	userID, err = a.cycleRefreshTokens(c)
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("user_id", userID)
	c.Next()
}

// cycles refresh tokens and returns the user associated with the token. Returns err if token not valid
func (a *Authenticator) cycleRefreshTokens(c *gin.Context) (uint, error) {
	// ID token is expired, check the refresh token
	refreshTokenString, err := c.Cookie("refresh_token")
	if err != nil {
		return 0, errors.New("refresh token could not be found")
	}
	refreshToken, err := a.tokenService.ValidateRefreshToken(refreshTokenString)
	if err != nil || refreshToken == nil {
		log.Print(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return 0, nil
	}
	//cycle refresh tokens
	rID := refreshToken.ID
	uid := refreshToken.UserID
	tokenPair, err := a.tokenService.GetNewTokenPair(uid, rID.String())
	if err != nil {
		return 0, errors.New("token pair could not be generated: " + err.Error())
	}
	a.updateCookies(c, tokenPair)
	return tokenPair.UserID, nil
}

func (a *Authenticator) updateCookies(c *gin.Context, tokenPair *models.TokenPair) {
	c.SetCookie("id_token", tokenPair.IDToken.SS, int(a.tokenService.IDTokenLifetime), "/", a.host, false, true)
	c.SetCookie("refresh_token", tokenPair.RefreshToken.SS, int(a.tokenService.RefreshTokenLifetime), "/", a.host, false, true)
}

func (a *Authenticator) resetAdminCache() {
	if a.writing {
		return
	}
	off, err := a.userService.FindOfficers()
	if err != nil {
		log.Fatal("could not load users: " + err.Error())
	}
	//fugly data condom
	if a.writing {
		return
	}
	a.cacheMutex.Lock()
	defer a.cacheMutex.Unlock()
	a.adminIDs = make(map[uint]bool)
	for _, u := range off {
		a.adminIDs[u.ID] = true
	}
}
