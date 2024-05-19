package routes

import (
	"ClubTennis/config"
	"ClubTennis/controllers"
	"ClubTennis/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func SetAuthRoutes(c config.Config, engine *gin.Engine, db *gorm.DB, client *redis.Client) {
	authGroup := engine.Group("/auth")
	var authCtrl controllers.AuthController = *controllers.NewAuthController(c, services.NewUserService(db))

	{
		authGroup.Use( /*middleware*/ )

		authGroup.POST("/login", authCtrl.Login)
		authGroup.GET("/callback", authCtrl.Callback)
	}
}
