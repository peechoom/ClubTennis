package routes

import (
	"ClubTennis/controllers"
	"ClubTennis/services"

	"github.com/gin-gonic/gin"
)

// sets up the routes for all endpoints that are meant to be used by admins only via /admin/ grouping
func SetAdminRoutes(engine *gin.Engine, s *services.ServiceContainer) {
	adminGroup := engine.Group("/admin")
	// var matchCtrl controllers.MatchController = *controllers.NewMatchController(db)
	var userCtrl controllers.UserController = *controllers.NewUserController(s.UserService, s.MatchService)
	// var auth middleware.Authenticator = *middleware.NewAuthenticator(s.TokenService, s.UserService, os.Getenv("SERVER_HOST"))

	{
		//TODO make admin accounts saveable to a file or sumn
		adminGroup.Use( /*auth.AuthenticateAdmin*/ )

		//admin webpage handlers
		adminGroup.GET("/", controllers.AdminHomeHandler)
		adminGroup.GET("/index.html", controllers.AdminHomeHandler)

		adminGroup.GET("/editmembers", controllers.EditMembersHandler)
		adminGroup.GET("/editmembers.html", controllers.EditMembersHandler)

		//API handlers
		//for matches
		//...

		//for users
		adminGroup.POST("/members", userCtrl.CreateMember)
		adminGroup.PUT("/members/:id", userCtrl.EditMember)
		adminGroup.DELETE("/members/:id", userCtrl.DeleteMember)
	}
}
