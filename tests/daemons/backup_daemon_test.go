package daemons_test

import (
	"ClubTennis/daemons"
	"ClubTennis/initializers"
	"ClubTennis/models"
	"ClubTennis/services"
	"testing"
	"time"
)

func TestAutoBackupDaemon(t *testing.T) {
	db := initializers.GetTestDatabase()
	if db == nil {
		panic("error in setup!")
	}
	db.Exec("DROP SCHEMA " + initializers.TestDBName + ";")

	db = initializers.GetTestDatabase()
	if db == nil {
		panic("error in setup!")
	}
	err := db.AutoMigrate(models.User{}, models.Match{})
	if err != nil {
		panic(err)
	}

	es := services.NewEmailService("/home/alec/go/src/ClubTennis/templates", "test@test.com", "")
	us := services.NewUserService(db)

	var users []*models.User
	u, _ := models.NewUser("poopie", "unc", "Fart", "Stinky", "poopie@gmail.com", models.MENS_LADDER)
	users = append(users, u)
	u, _ = models.NewUser("pdiddy", "unc", "Puff", "Diddy", "pdiddle@unc.edu", models.MENS_LADDER)
	users = append(users, u)
	u, _ = models.NewUser("jembitch", "sixers", "Joel", "Embitch", "jembitch@sixers.com", models.WOMENS_LADDER)
	users = append(users, u)
	u, _ = models.NewUser("sbarnes", "raptors", "Scottie", "Barnes", "shbarnes@sixers.com", models.MENS_LADDER)
	users = append(users, u)
	us.Save(users...)

	go daemons.AutoBackupDaemon(time.Second*2, false, es, us)
	time.Sleep(time.Second * 3)
	//u gotta put ur own email creds in it to send

}
