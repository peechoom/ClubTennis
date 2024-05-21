package routes

import (
	"ClubTennis/controllers"
	"ClubTennis/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetAuthRoutes(engine *gin.Engine, db *gorm.DB) {
	authGroup := engine.Group("/auth")
	var authCtrl controllers.AuthController = *controllers.NewAuthController(services.NewUserService(db))

	{
		authGroup.Use( /*middleware*/ )

		authGroup.POST("/login", authCtrl.Login)
		authGroup.GET("/login", authCtrl.Login)
		authGroup.POST("/callback", authCtrl.Callback)
	}
}
