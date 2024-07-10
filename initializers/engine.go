package initializers

import (
	"ClubTennis/daemons"
	"ClubTennis/routes"
	"ClubTennis/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetEngine() *gin.Engine {
	e := gin.New()
	e.Use(gin.Recovery())
	//e.Use(gin.Logger())

	e = e.Delims("{[{", "}]}")
	e.LoadHTMLGlob("templates/*.html")

	s := services.SetupServices(GetDatabase(), "templates")

	setRoutings(e, s)
	startDaemons(s)

	//probably best to leave this off. might be needed for caddy?
	e.SetTrustedProxies(nil)
	return e
}

func GetTestEngine() *gin.Engine {
	e := gin.Default()

	//TODO make absolute path an env variable that we fetch
	e = e.Delims("{[{", "}]}")
	e.LoadHTMLGlob("/home/alec/go/src/ClubTennis/templates/*.html")

	s := services.SetupServices(GetTestDatabase(), "templates")
	setRoutings(e, s)

	return e
}

func setRoutings(e *gin.Engine, s *services.ServiceContainer) {
	// ping endpoint returns unix time
	e.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusAccepted, "%d", time.Now().Unix())
	})
	routes.SetAuthRoutes(e, s)
	routes.SetAdminRoutes(e, s)
	routes.SetClubRoutes(e, s)
	routes.SetPublicRoutes(e, s)
}

func startDaemons(s *services.ServiceContainer) {
	go daemons.MatchWarningDaemon(daemons.WarningDefaultFrequency, true, s.MatchService, s.EmailService)
	go daemons.MatchExpiredDaemon(daemons.ExpiredDefaultFrequency, true, s.MatchService, s.EmailService)
	go daemons.MatchDeletionDaemon(daemons.DeletionDefaultFrequency, true, s.MatchService)
	go daemons.AutoBackupDaemon(daemons.AutoBackupDefaultFrequency, true, s.EmailService, s.UserService)
}
