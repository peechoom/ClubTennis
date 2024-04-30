package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UnityID     string   //ncsu unity id or skema id, should be unique
	Affiliation string   //ncsu.edu or skema.edu
	FirstName   string   //users first name
	LastName    string   //users last name
	Rank        uint     `gorm:"uniqueIndex,sort:asc"` //users rank in the ladder
	Wins        uint     //how many wins the player has
	Losses      uint     //how many losses the player has
	Matches     []*Match //list of matches the player is involved in
}

/*
*

	create a new user
*/
func NewUser(UnityID string, Affiliation string, FirstName string, LastName string) *User {
	if UnityID == "" || Affiliation == "" || FirstName == "" || LastName == "" {
		return nil
	}
	u := new(User)

	u.UnityID = UnityID
	u.Affiliation = Affiliation
	u.FirstName = FirstName
	u.LastName = LastName
	return u
}
