package routes

import (
	"ClubTennis/controllers"
	"ClubTennis/middleware"
	"ClubTennis/services"

	"github.com/gin-gonic/gin"
)

func SetAuthRoutes(engine *gin.Engine, s *services.ServiceContainer) {
	authGroup := engine.Group("/auth")
	var authCtrl *controllers.AuthController = controllers.NewAuthController(s.UserService, s.TokenService)

	{
		authGroup.Use(middleware.DoNotCache)

		authGroup.POST("/login", authCtrl.Login)
		authGroup.GET("/login", authCtrl.Login)
		authGroup.GET("/callback", authCtrl.Callback)
		authGroup.GET("/me", authCtrl.Me)
		authGroup.GET("/logout", authCtrl.Logout)
		authGroup.POST("/logout", authCtrl.Logout)
	}
}
