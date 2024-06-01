package routes

import (
	"ClubTennis/controllers"
	"ClubTennis/services"

	"github.com/gin-gonic/gin"
)

// sets up the routings for endpoints that anyone can access
func SetPublicRoutes(engine *gin.Engine, s *services.ServiceContainer) {
	r := engine
	pubCtrl := controllers.NewPublicController(s.PublicService)
	{
		r.GET("/", controllers.HomeHandler)
		r.GET("/index.html", controllers.HomeHandler)

		r.GET("/about", controllers.AboutHandler)
		r.GET("/about.html", controllers.AboutHandler)

		r.GET("/signin", controllers.SigninHandler)
		r.GET("/signin.html", controllers.SigninHandler)

		r.GET("/slides", pubCtrl.GetSlideshow)
		r.GET("/welcome", pubCtrl.GetWelcome)

		r.Static("/static", "./static")
		r.StaticFile("/favicon.ico", "./favicon.ico")
	}
}
