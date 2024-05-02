package models

import (
	"errors"
	"regexp"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UnityID     string   `gorm:"index:,unique,sort:desc,type:btree,length:50"` //ncsu unity id or skema id, should be unique
	Affiliation string   //ncsu.edu or skema.edu
	FirstName   string   //users first name
	LastName    string   //users last name
	Email       string   //users e-mail address
	Rank        uint     `gorm:"index:,sort:desc"` //users rank in the ladder, should be unique
	Wins        int      //how many wins the player has
	Losses      int      //how many losses the player has
	Matches     []*Match `gorm:"many2many:user_matches;"` //list of matches the player is involved in
	isOfficer   bool     //whether or not this user is an officer
}

// thanks openAI
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

/*
create a new user
*/
func NewUser(UnityID string, Affiliation string, FirstName string, LastName string, Email string) (*User, error) {
	if UnityID == "" || Affiliation == "" || FirstName == "" || LastName == "" {
		return nil, errors.New("some fields not given")
	}
	if !isValidEmail(Email) {
		return nil, errors.New("email not valid")
	}

	u := new(User)

	u.UnityID = UnityID
	u.Affiliation = Affiliation
	u.FirstName = FirstName
	u.LastName = LastName
	u.Email = Email
	u.isOfficer = false
	return u, nil
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

/*
create a new officer
*/
func NewOfficer(UnityID string, Affiliation string, FirstName string, LastName string, Email string) (*User, error) {
	u, err := NewUser(UnityID, Affiliation, FirstName, LastName, Email)
	if err != nil {
		return nil, err
	}
	u.isOfficer = true
	return u, nil
}

/*
promotes user to officer
*/
func (u *User) SetOfficer(isOfficer bool) {
	u.isOfficer = isOfficer
}

/*
demotes officer
*/
func (u *User) IsOfficer() bool {
	return u.isOfficer
}
