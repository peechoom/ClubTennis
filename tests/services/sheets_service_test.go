package services_test

import (
	"ClubTennis/models"
	"ClubTennis/services"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUsersToSheet(t *testing.T) {
	var users []models.User
	u, _ := models.NewUser("poopie", "unc", "Fart", "Stinky", "poopie@gmail.com", models.MENS_LADDER)
	users = append(users, *u)
	u, _ = models.NewUser("pdiddy", "unc", "Puff", "Diddy", "pdiddle@unc.edu", models.MENS_LADDER)
	users = append(users, *u)
	u, _ = models.NewUser("jembitch", "sixers", "Joel", "Embitch", "jembitch@sixers.com", models.WOMENS_LADDER)
	users = append(users, *u)
	u, _ = models.NewUser("sbarnes", "raptors", "Scottie", "Barnes", "shbarnes@sixers.com", models.MENS_LADDER)
	users = append(users, *u)

	t.Log(services.UsersToSheet(users))
}

func TestSheetToUsers(t *testing.T) {
	var users []models.User
	u, _ := models.NewUser("poopie", "unc", "Fart", "Stinky", "poopie@gmail.com", models.MENS_LADDER)
	users = append(users, *u)
	u, _ = models.NewUser("pdiddy", "unc", "Puff", "Diddy", "pdiddle@unc.edu", models.MENS_LADDER)
	users = append(users, *u)
	u, _ = models.NewUser("jembitch", "sixers", "Joel", "Embitch", "jembitch@sixers.com", models.WOMENS_LADDER)
	users = append(users, *u)
	u, _ = models.NewUser("sbarnes", "raptors", "Scottie", "Barnes", "shbarnes@sixers.com", models.MENS_LADDER)
	users = append(users, *u)
	filename, err := services.UsersToSheet(users)
	require.NoError(t, err)

	loaded, err := services.SheetToUsers(filename)
	require.NoError(t, err)
	for _, u := range users {
		require.Contains(t, loaded, u)
	}

}
