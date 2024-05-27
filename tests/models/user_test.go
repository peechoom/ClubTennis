package models_test

import (
	"ClubTennis/models"
	"testing"

	"github.com/stretchr/testify/require"
)

// test basic user functionality
func TestNewUser(t *testing.T) {
	user, err := models.NewUser("abcde6", "ncsu", "John", "Doe", "asdasdfsfd4@ncsu.edu")

	require.NoError(t, err)
	require.Equal(t, user.FirstName, "John")
	require.Equal(t, user.LastName, "Doe")
	require.Equal(t, user.Affiliation, "ncsu")
	require.Equal(t, user.UnityID, "abcde6")
	require.Equal(t, user.Email, "asdasdfsfd4@ncsu.edu")
	require.False(t, user.IsOfficer)

	user.SetOfficer(true)
	require.True(t, user.IsOfficer)

	o, err := models.NewOfficer("abcde6", "skema", "John", "Doe", "asfdfsfdesfsd.abdfsfsfsfdsfd@skema.edu")

	require.NoError(t, err)
	require.Equal(t, o.FirstName, "John")
	require.Equal(t, o.LastName, "Doe")
	require.Equal(t, o.Affiliation, "skema")
	require.Equal(t, o.UnityID, "abcde6")
	require.Equal(t, o.Email, "asfdfsfdesfsd.abdfsfsfsfdsfd@skema.edu")

	require.True(t, o.IsOfficer)

	o.SetOfficer(false)
	require.False(t, o.IsOfficer)
}

// test new user with bad email
func TestNewUserBadEmail(t *testing.T) {
	user, err := models.NewUser("abcdef", "ncsu", "John", "Doe", "notvalid@notanEmail")

	require.Nil(t, user)
	require.Error(t, err)

	user, err = models.NewUser("abcdef", "ncsu", "John", "Doe", "notvalid@notanEmail.")

	require.Nil(t, user)
	require.Error(t, err)

	user, err = models.NewUser("abcdef", "ncsu", "John", "Doe", "notanEmail.com")

	require.Nil(t, user)
	require.Error(t, err)

	user, err = models.NewUser("abcdef", "ncsu", "John", "Doe", "@notAnEmail")

	require.Nil(t, user)
	require.Error(t, err)

	user, err = models.NewUser("abcdef", "ncsu", "John", "Doe", "notAnEmail@")

	require.Nil(t, user)
	require.Error(t, err)
}
