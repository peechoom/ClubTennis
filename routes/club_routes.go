package routes

import (
	"ClubTennis/controllers"
	"ClubTennis/middleware"
	"ClubTennis/services"
	"os"

	"github.com/gin-gonic/gin"
)

// sets up the routes for any club member via the /club/ grouping
func SetClubRoutes(engine *gin.Engine, s *services.ServiceContainer) {
	clubGroup := engine.Group("/club")
	var matchCtrl *controllers.MatchController = controllers.NewMatchController(s.MatchService, s.UserService, s.EmailService)
	var userCtrl *controllers.UserController = controllers.NewUserController(s.UserService, s.MatchService)
	var annCtrl *controllers.AnnouncementController = controllers.NewAnnouncementController(s.AnnouncementService, s.EmailService, s.UserService, s.ImageService)
	var auth *middleware.Authenticator = middleware.NewAuthenticator(s.TokenService, s.UserService, os.Getenv("SERVER_HOST"))
	var pageCtrl *controllers.PageController = controllers.NewPageController(s.SnippetService)
	var rulesCtrl *controllers.RulesController = controllers.NewRulesController(s.ImageService, s.SnippetService, os.Getenv("SERVER_HOST"))
	{
		clubGroup.Use(auth.AuthenticateMember)

		//club webpage handlers
		clubGroup.GET("/challenge", controllers.ChallengeHandler)
		clubGroup.GET("/challenge.html", controllers.ChallengeHandler)

		clubGroup.GET("/announcements", controllers.AnnouncementsHandler)
		clubGroup.GET("/announcements.html", controllers.AnnouncementsHandler)

		clubGroup.GET("/", controllers.ClubHomeHandler)
		clubGroup.GET("/index.html", controllers.ClubHomeHandler)

		clubGroup.GET("/challengerules", pageCtrl.ChallengeRulesHandler)
		clubGroup.GET("/challengerules.html", pageCtrl.ChallengeRulesHandler)

		clubGroup.GET("/ladderrules", pageCtrl.LadderRulesHandler)
		clubGroup.GET("/ladderrules.html", pageCtrl.LadderRulesHandler)

		// API handlers
		// for matches
		clubGroup.POST("/matches", matchCtrl.Challenge)
		clubGroup.GET("/matches/recent", matchCtrl.GetRecentMatches)
		clubGroup.PATCH("/matches/:id", matchCtrl.SubmitScore)
		clubGroup.GET("/matches/:id", matchCtrl.GetMatchByID)
		clubGroup.GET("/matches", matchCtrl.GetMatch)

		//for users
		clubGroup.GET("/members/:id", userCtrl.GetMemberByID)
		clubGroup.GET("/members", userCtrl.GetAllMembers)
		clubGroup.GET("/cutoff/:ladder", userCtrl.GetCutoff)

		// for announcements
		clubGroup.GET("/announcements/:page", annCtrl.GetAnnouncementPage)
		// for snippets
		clubGroup.GET("/ladderrulessnippet", rulesCtrl.GetLadderRules)
		clubGroup.GET("/challengerulessnippet", rulesCtrl.GetChallengeRules)

	}
}
