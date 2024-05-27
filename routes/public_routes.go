package routes

import (
	"ClubTennis/controllers"

	"github.com/gin-gonic/gin"
)

// sets up the routings for endpoints that anyone can access
func SetPublicRoutes(engine *gin.Engine) {
	r := engine

	{
		r.GET("/", controllers.HomeHandler)
		r.GET("/index.html", controllers.HomeHandler)

		r.GET("/about", controllers.AboutHandler)
		r.GET("/about.html", controllers.AboutHandler)

		r.GET("/contact", controllers.ContactHandler)
		r.GET("/contact.html", controllers.ContactHandler)

		r.GET("/signin", controllers.SigninHandler)
		r.GET("/signin.html", controllers.SigninHandler)

		r.Static("/static", "./static")
		r.StaticFile("/favicon.ico", "/static/favicon.ico")
	}
}
