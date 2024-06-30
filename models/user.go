package models

import (
	"errors"
	"regexp"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UnityID         string   `gorm:"index:,unique,sort:desc,type:btree,length:50"` //ncsu unity id or skema id, should be unique
	Affiliation     string   //ncsu.edu or skema.edu
	FirstName       string   //users first name
	LastName        string   //users last name
	SigninEmail     string   `gorm:"index:,unique,sort:desc,type:btree,length:100"` //users e-mail address
	ContactEmail    string   //users email address for sending emails to (is usually the same as signin email)
	Rank            uint     `gorm:"index:,sort:desc"` //users rank in the ladder, should be unique
	Wins            int      //how many wins the player has
	Losses          int      //how many losses the player has
	Matches         []*Match `gorm:"constraint:OnDelete:CASCADE;many2many:user_matches"` //list of matches the player is involved in
	IsOfficer       bool     //whether or not this user is an officer
	IsChallengeable bool     `gorm:"-:all"` //can this user be challenged? not stored
	Ladder          string   // what ladder this user plays in -> 'M' for mens, 'W' for womens
}

const MENS_LADDER string = "M"
const WOMENS_LADDER string = "W"

// thanks openAI
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

/*
create a new user
*/
func NewUser(UnityID string, Affiliation string, FirstName string, LastName string, Email string, Ladder string) (*User, error) {
	if UnityID == "" || Affiliation == "" || FirstName == "" || LastName == "" {
		return nil, errors.New("some fields not given")
	}
	if !isValidEmail(Email) {
		return nil, errors.New("email not valid")
	}
	if len(UnityID) > 50 {
		return nil, errors.New("unity id too long")
	}
	u := new(User)

	u.UnityID = UnityID
	u.Affiliation = Affiliation
	u.FirstName = FirstName
	u.LastName = LastName
	u.SigninEmail = Email
	u.IsOfficer = false
	u.Ladder = Ladder
	u.ContactEmail = u.SigninEmail
	return u, nil
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email) && len(email) < 100
}

/*
create a new officer
*/
func NewOfficer(UnityID string, Affiliation string, FirstName string, LastName string, Email string, Ladder string) (*User, error) {
	u, err := NewUser(UnityID, Affiliation, FirstName, LastName, Email, Ladder)
	if err != nil {
		return nil, err
	}
	u.IsOfficer = true
	return u, nil
}

/*
sets officer status
*/
func (u *User) SetOfficer(isOfficer bool) {
	u.IsOfficer = isOfficer
}

// edits the fields of this user to be the fields of the new user. cannot edit ID because it is a primary key.
// fields that are 0 length in the new user are assumed to be unchanged, wins/losses will always be updated tho. This should really only be called by
// officers. This function is for critical stuff like unity id, ranking, etc. that members themselves shouldnt control.
func (u *User) EditUser(nu *User) {
	if len(nu.UnityID) != 0 {
		u.UnityID = nu.UnityID
	}
	if len(nu.Affiliation) != 0 {
		u.Affiliation = nu.Affiliation
	}
	if len(nu.FirstName) != 0 {
		u.FirstName = nu.FirstName
	}
	if len(nu.LastName) != 0 {
		u.LastName = nu.LastName
	}
	if len(nu.SigninEmail) != 0 && isValidEmail(nu.SigninEmail) {
		u.SigninEmail = nu.SigninEmail
	}
	if nu.Rank != 0 {
		u.Rank = nu.Rank
	}
	if nu.Ladder != "" {
		u.Ladder = nu.Ladder
	}
	if nu.ContactEmail != "" {
		u.ContactEmail = nu.ContactEmail
	}
	u.Wins = nu.Wins
	u.Losses = nu.Losses
	u.SetOfficer(nu.IsOfficer)
}
