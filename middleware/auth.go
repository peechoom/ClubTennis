package middleware

import (
	"ClubTennis/models"
	"ClubTennis/repositories"
	"ClubTennis/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Authenticator struct {
	tokenService *services.TokenService
	userService  *services.UserService
	host         string
}

func NewAuthenticator(tokenRepo *repositories.TokenRepository, db *gorm.DB, host string) *Authenticator {
	return &Authenticator{
		tokenService: services.DefaultTokenService(tokenRepo),
		userService:  services.NewUserService(db),
		host:         host,
	}
}

func (a *Authenticator) AuthenticateMember(c *gin.Context) {
	var err error
	idTokenString, err := c.Cookie("id_token")
	if err != nil {
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
	if err != nil {
		c.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var userID uint
	userID, err = a.tokenService.ValidateIDToken(idTokenString)
	//TODO check db to verify admin status
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

// cycles refresh tokens and returns the user associated with the token. Returns err if token not valid
func (a *Authenticator) cycleRefreshTokens(c *gin.Context) (uint, error) {
	// ID token is expired, check the refresh token
	refreshTokenString, err := c.Cookie("refresh_token")
	if err != nil {
		return 0, errors.New("you are not authorized to access this page")
	}
	refreshToken, err := a.tokenService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		// refresh token expired, user needs to sign in again
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	//cycle refresh tokens
	tokenPair, err := a.tokenService.GetNewTokenPair(refreshToken.UserID, refreshToken.SS)
	if err != nil {
		return 0, errors.New("you are not authorized to access this page")
	}
	a.updateCookies(c, tokenPair)
	return tokenPair.UserID, nil
}

func (a *Authenticator) updateCookies(c *gin.Context, tokenPair *models.TokenPair) {
	c.SetCookie("id_token", tokenPair.IDToken.SS, int(a.tokenService.IDTokenLifetime), "/", a.host, false, true)
	c.SetCookie("refresh_token", tokenPair.RefreshToken.SS, int(a.tokenService.RefreshTokenLifetime), "/", a.host, false, true)
}
