package routes

import (
	"ClubTennis/controllers"
	"ClubTennis/middleware"
	"ClubTennis/services"
	"os"

	"github.com/gin-gonic/gin"
)

// sets up the routes for all endpoints that are meant to be used by admins only via /admin/ grouping
func SetAdminRoutes(engine *gin.Engine, s *services.ServiceContainer) {
	adminGroup := engine.Group("/admin")
	// var matchCtrl controllers.MatchController = *controllers.NewMatchController(db)
	var userCtrl *controllers.UserController = controllers.NewUserController(s.UserService, s.MatchService)
	var auth middleware.Authenticator = *middleware.NewAuthenticator(s.TokenService, s.UserService, os.Getenv("SERVER_HOST"))
	var annCtrl *controllers.AnnouncementController = controllers.NewAnnouncementController(s.AnnouncementService, s.EmailService, s.UserService, s.ImageService)
	var pubCtrl *controllers.PublicController = controllers.NewPublicController(s.PublicService, s.ImageService)
	var backCtrl *controllers.BackupController = controllers.NewBackupController(s.UserService)
	{
		adminGroup.Use(auth.AuthenticateAdmin)

		//admin webpage handlers
		adminGroup.GET("/", controllers.AdminHomeHandler)
		adminGroup.GET("/index.html", controllers.AdminHomeHandler)

		adminGroup.GET("/editmembers", controllers.EditMembersHandler)
		adminGroup.GET("/editmembers.html", controllers.EditMembersHandler)

		adminGroup.GET("/sendannouncements", controllers.SendAnnouncementsHandler)
		adminGroup.GET("/sendannouncements.html", controllers.SendAnnouncementsHandler)

		adminGroup.GET("/editpublicpages", controllers.EditPublicPagesHandler)
		adminGroup.GET("/editpublicpages.html", controllers.EditPublicPagesHandler)

		//API handlers
		//for matches
		//...

		//for users
		adminGroup.POST("/members", userCtrl.CreateMember)
		adminGroup.PUT("/members/:id", userCtrl.EditMember)
		adminGroup.DELETE("/members/:id", userCtrl.DeleteMember)
		backupGroup := adminGroup.Group("/backups")
		{
			backupGroup.GET("/users", backCtrl.GetBackupSpreadsheet)
			backupGroup.POST("/users", backCtrl.LoadBackupSpreadsheet)
		}

		// for announcements
		adminGroup.POST("/announcements", annCtrl.SubmitPost)
		adminGroup.DELETE("/announcements/:id", annCtrl.DeleteAnnouncement)

		//for public facing things
		adminGroup.POST("/slides/:slideNum", pubCtrl.PostSlides)
		adminGroup.PUT("/welcome", pubCtrl.PutWelcome)
	}
}
