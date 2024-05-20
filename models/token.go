package models

import "github.com/google/uuid"

// revlovling refresh token that the user can use in order to get a new ID token from the system.
// saved to in-memory db
type RefreshToken struct {
	ID     uuid.UUID `json:"-"`
	UserID uint      `json:"-"`
	SS     string    `json:"refreshToken"`
}

// ID token that the user uses to identify itself. not saved.
type IDToken struct {
	SS string `json:"idToken"`
}

// Simple touple containing an ID Token and a Refresh Token. Is sent by the server
type TokenPair struct {
	IDToken
	RefreshToken
}
