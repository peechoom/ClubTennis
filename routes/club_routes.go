package routes

import (
	"ClubTennis/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// sets up the routes for any club member via the /club/ grouping
func SetClubRoutes(engine *gin.Engine, db *gorm.DB) {
	clubGroup := engine.Group("/club")
	var matchCtrl controllers.MatchController = *controllers.NewMatchController(db)
	var userCtrl controllers.UserController = *controllers.NewUserController(db)
	{
		clubGroup.Use( /*club auth middleware*/ )

		//club webpage handlers
		clubGroup.GET("/challenge", controllers.ChallengeHandler)
		clubGroup.GET("/challenge.html", controllers.ChallengeHandler)

		clubGroup.GET("/announcements", controllers.AnnouncementsHandler)
		clubGroup.GET("/announcements.html", controllers.AnnouncementsHandler)

		clubGroup.GET("/", controllers.ClubHomeHandler)
		clubGroup.GET("/index.html", controllers.ClubHomeHandler)

		// API handlers
		// for matches
		clubGroup.POST("/matches", matchCtrl.Challenge)
		clubGroup.PATCH("/matches/:id", matchCtrl.SubmitScore)
		clubGroup.GET("/matches/:id", matchCtrl.GetMatchByID)
		clubGroup.GET("/matches", matchCtrl.GetMatch)

		//for users
		clubGroup.GET("/members/:id", userCtrl.GetMemberByID)
		clubGroup.GET("/members", userCtrl.GetAllMembers)

	}
}
