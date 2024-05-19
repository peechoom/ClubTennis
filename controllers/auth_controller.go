package controllers

import (
	"ClubTennis/config"
	"ClubTennis/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mazen160/go-random"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type AuthController struct {
	googleOauthConfig *oauth2.Config
	userService       *services.UserService
	tokenService      *services.TokenService
	host              string //this host
	stateString       string
}

func NewAuthController(c config.Config, userService *services.UserService) *AuthController {
	statestr, err := random.String(64)
	if err != nil {
		return nil
	}
	return &AuthController{
		googleOauthConfig: &oauth2.Config{
			RedirectURL:  c.GoogleOauth.RedirectURL,
			ClientID:     c.GoogleOauth.ClientID,
			ClientSecret: c.GoogleOauth.ClientSecret,
			Endpoint:     google.Endpoint,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		},
		stateString: statestr,
		userService: userService,
		host:        c.Server.Host,
	}
}

/*
	POST /auth/login

users POST to here when they want to log in
*/
func (a *AuthController) Login(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, a.googleOauthConfig.AuthCodeURL(a.stateString))
}

//TODO add logout

/*
	GET /auth/callback

google will GET this and redirect user to it when login success
*/
func (a *AuthController) Callback(c *gin.Context) {
	state := c.Query("state")
	if state != a.stateString {
		c.Error(errors.New("state string does not match"))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := c.Query("code")
	email, err := getEmailFromGoogleToken(c, code, a.googleOauthConfig)
	if err != nil || email == "" {
		c.Error(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	user, err := a.userService.FindByEmail(email)
	if err != nil {
		c.Error(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	c.Set("user_id", user.ID)

	tokenPair, err := a.tokenService.GetNewTokenPair(user.ID, "")
	if err != nil {
		c.Error(err)
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	//TODO update this to use the config host
	c.SetCookie("id_token", tokenPair.IDToken.SS, int(a.tokenService.IDTokenLifetime), "/", a.host, false, true)
	c.SetCookie("refresh_token", tokenPair.RefreshToken.SS, int(a.tokenService.RefreshTokenLifetime), "/", a.host, false, true)

	if user.IsOfficer() {
		c.Redirect(http.StatusPermanentRedirect, "/admin/")
		return
	} else {
		c.Redirect(http.StatusPermanentRedirect, "/club/")
		return
	}
}

// UTILITY FUNCTIONS

// gets email address from a google token
func getEmailFromGoogleToken(c *gin.Context, codeString string, config *oauth2.Config) (string, error) {
	token, err := config.Exchange(c.Request.Context(), codeString)
	if err != nil {
		return "", err
	}
	if !token.Valid() {
		return "", errors.New("oauth token is expired")
	}

	client := config.Client(c.Request.Context(), token)
	if client == nil {
		return "", errors.New("error getting client from token")
	}
	oauth2Service, err := oauth2api.NewService(c.Request.Context(), option.WithHTTPClient(client))
	if err != nil {
		return "", err
	}
	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return "", err
	}

	return userInfo.Email, nil
}
