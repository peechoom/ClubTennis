package initializers

import (
	"ClubTennis/routes"
	"ClubTennis/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetEngine() *gin.Engine {
	e := gin.Default()

	//TODO make absolute path an env variable that we fetch
	e = e.Delims("{[{", "}]}")
	e.LoadHTMLGlob("templates/*.html")
	db := GetDatabase()
	setRoutings(e, db)

	return e
}

func GetTestEngine() *gin.Engine {
	e := gin.Default()

	//TODO make absolute path an env variable that we fetch
	e = e.Delims("{[{", "}]}")
	e.LoadHTMLGlob("/home/alec/go/src/ClubTennis/templates/*.html")

	db := GetTestDatabase()
	setRoutings(e, db)

	return e
}

func setRoutings(e *gin.Engine, db *gorm.DB) {
	s := services.SetupServices(db, "templates")
	routes.SetAuthRoutes(e, s)
	routes.SetAdminRoutes(e, s)
	routes.SetClubRoutes(e, s)
	routes.SetPublicRoutes(e, s)
}
