package routes

import (
	"ClubTennis/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// sets up the routes for all endpoints that are meant to be used by admins only via /admin/ grouping
func SetAdminRoutes(engine *gin.Engine, db *gorm.DB) {
	adminGroup := engine.Group("/admin")
	// var matchCtrl controllers.MatchController = *controllers.NewMatchController(db)
	var userCtrl controllers.UserController = *controllers.NewUserController(db)

	{
		adminGroup.Use( /*admin middleware for authorization*/ )

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
