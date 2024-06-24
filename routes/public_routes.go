package routes

import (
	"ClubTennis/controllers"
	"ClubTennis/middleware"
	"ClubTennis/services"

	"github.com/gin-gonic/gin"
)

// sets up the routings for endpoints that anyone can access
func SetPublicRoutes(engine *gin.Engine, s *services.ServiceContainer) {
	r := engine
	pubCtrl := controllers.NewPublicController(s.PublicService, s.ImageService)
	{

		r.GET("/", controllers.HomeHandler)
		r.GET("/index.html", controllers.HomeHandler)

		r.GET("/about", controllers.AboutHandler)
		r.GET("/about.html", controllers.AboutHandler)

		r.GET("/signin", controllers.SigninHandler)
		r.GET("/signin.html", controllers.SigninHandler)

		r.GET("/slides", pubCtrl.GetSlideshow)
		r.GET("/welcome", pubCtrl.GetWelcome)

		r.NoRoute(controllers.ErrorHandler)
		r.GET("/error", controllers.ErrorHandler)
		r.GET("/error.html", controllers.ErrorHandler)

		staticGroup := r.Group("/static")
		{
			staticGroup.Use(middleware.GetCors())
			staticGroup.Static("/", "./static")
			r.StaticFile("/favicon.ico", "./favicon.ico")
		}

		//for getting images.
		r.GET("/images/:filename", pubCtrl.GetImage)
	}
}
