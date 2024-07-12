package controllers

import (
	"ClubTennis/models"
	"ClubTennis/services"
	"errors"
	"log"
	"net/http"
	"os"

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
	host              string // this host
	stateString       string
}

func NewAuthController(userService *services.UserService, tokenservice *services.TokenService) *AuthController {
	statestr, err := random.String(64)
	if err != nil {
		return nil
	}
	return &AuthController{
		googleOauthConfig: &oauth2.Config{
			RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
			ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
			Endpoint:     google.Endpoint,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		},
		stateString:  statestr,
		userService:  userService,
		tokenService: tokenservice,
		host:         os.Getenv("SERVER_HOST"),
	}
}

func (a *AuthController) Login(c *gin.Context) {

	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, a.googleOauthConfig.AuthCodeURL(a.stateString))
		return
	}
	refreshToken, err := a.tokenService.ValidateRefreshToken(cookie)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, a.googleOauthConfig.AuthCodeURL(a.stateString))
		return
	}
	tokens, err := a.tokenService.GetNewTokenPair(refreshToken.UserID, refreshToken.ID.String())
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, a.googleOauthConfig.AuthCodeURL(a.stateString))
		return
	}
	_, err = a.userService.FindByID(tokens.UserID)
	if err != nil {
		c.Redirect(http.StatusPermanentRedirect, "/")
		return
	}
	setCookies(c, tokens, int(a.tokenService.IDTokenLifetime), int(a.tokenService.RefreshTokenLifetime), a.host)

	if tokens.UserID > 0 {
		c.Redirect(http.StatusPermanentRedirect, "/club/")
		return
	} else {
		c.Redirect(http.StatusPermanentRedirect, "/admin/")
		return
	}
}

func (a *AuthController) Logout(c *gin.Context) {
	host := a.host // Use the host from the controller

	tokenStr, err := c.Cookie("refresh_token")
	if err != nil {
		c.Redirect(http.StatusPermanentRedirect, "/")
		return
	}
	token, _ := a.tokenService.ValidateRefreshToken(tokenStr)

	a.tokenService.DeleteRefreshToken(token.UserID, token.ID.String())

	// Log the cookies being set for debugging purposes
	log.Println("Clearing cookies")
	log.Println("Host:", host)

	// Set the cookies with a max age of -1 to delete them
	c.SetCookie("id_token", "", -1, "/", host, false, true)
	c.SetCookie("refresh_token", "", -1, "/", host, false, true)

	c.Redirect(http.StatusPermanentRedirect, "/")
}

func (a *AuthController) Callback(c *gin.Context) {
	state := c.Query("state")
	if state != a.stateString {
		c.Error(errors.New("state string does not match"))
		log.Print("state string does not match")
		c.Redirect(http.StatusTemporaryRedirect, "/error")
		return
	}

	code := c.Query("code")
	email, err := getEmailFromGoogleToken(c, code, a.googleOauthConfig)
	if err != nil || email == "" {
		c.Error(err)
		log.Print("no code")
		c.Redirect(http.StatusTemporaryRedirect, "/error")
		return
	}

	user, err := a.userService.FindBySigninEmail(email)
	if err != nil || user == nil {
		if email != os.Getenv("EMAIL_USERNAME") {
			c.Error(err)
			log.Print("no email")
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}
		user = &models.User{IsOfficer: true}
		user.ID = 0
	}
	c.Set("user_id", user.ID)

	tokenPair, err := a.tokenService.GetNewTokenPair(user.ID, "")
	if err != nil {
		c.Error(err)
		log.Print("couldn't make token pair")
		c.Redirect(http.StatusTemporaryRedirect, "/error")
		return
	}

	setCookies(c, tokenPair, int(a.tokenService.IDTokenLifetime), int(a.tokenService.RefreshTokenLifetime), a.host)

	// redirect root to admin page. micro optimization putting club first
	if user.ID > 0 {
		c.Redirect(http.StatusPermanentRedirect, "/club/")
		return
	} else {
		c.Redirect(http.StatusPermanentRedirect, "/admin/")
		return
	}
}

func (a *AuthController) Me(c *gin.Context) {
	ss, err := c.Cookie("id_token")
	if err != nil {
		c.Error(err)
		log.Print(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	uid, err := a.tokenService.ValidateIDToken(ss)
	if err != nil || uid == 0 {
		if err != nil {
			log.Print(err.Error())
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "id token not valid"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": uid})
}

// Utility functions

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

func setCookies(c *gin.Context, tokenPair *models.TokenPair, IDLifetime, RefreshLifetime int, host string) {
	c.SetCookie("id_token", tokenPair.IDToken.SS, IDLifetime, "/", host, false, true)
	c.SetCookie("refresh_token", tokenPair.RefreshToken.SS, RefreshLifetime, "/", host, false, true)

	// Log the cookies being set for debugging purposes
	log.Println("Setting cookies")
	log.Println("id_token:", tokenPair.IDToken.SS)
	log.Println("refresh_token:", tokenPair.RefreshToken.SS)
	log.Println("Host:", host)
}
