package repositories

import (
	"errors"
	"fmt"
	"time"
)

type TokenRepository struct {
	db *TTLRbTree
}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{db: NewTTLRbTree()}
}

// sets a users refresh token
func (repo *TokenRepository) SetRefreshToken(userID string, tokenID string, expiresIn time.Duration) {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	repo.db.Set(key, expiresIn)
	go repo.db.Clean()
}

// deletes the provided refresh token ID from the redis db.
// If the token was not found in the db, an error is returned, as the user
// has not provided a valid key
func (repo *TokenRepository) DeleteRefreshToken(userID string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	ok := repo.db.Del(key)
	go repo.db.Clean()

	if !ok {
		return errors.New("user not found")
	}
	return nil
}
