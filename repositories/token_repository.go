package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type TokenRepository struct {
	Redis *redis.Client
}

func NewTokenRepository(c *redis.Client) *TokenRepository {
	return &TokenRepository{Redis: c}
}

// sets a users refresh token
func (repo *TokenRepository) SetRefreshToken(userID string, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	err := repo.Redis.Set(key, 0, expiresIn).Err()
	if err != nil {
		return err
	}
	return nil
}

// deletes the provided refresh token ID from the redis db.
// If the token was not found in the db, an error is returned, as the user
// has not provided a valid key
func (repo *TokenRepository) DeleteRefreshToken(userID string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	result := repo.Redis.Del(key)

	if result.Err() != nil {
		return result.Err()
	}

	if result.Val() < 1 {
		return errors.New("invalid refresh token")
	}

	return nil
}

// deletes all of a users refresh tokens
func (repo *TokenRepository) DeleteUserRefreshTokens(userID string) error {
	pattern := fmt.Sprintf("%s*", userID)

	it := repo.Redis.Scan(0, pattern, 10).Iterator()
	var fails bool = false

	for it.Next() {
		err := repo.Redis.Del(it.Val()).Err()
		if err != nil {
			fails = true
		}
	}

	if it.Err() != nil {
		return it.Err()
	}

	if fails {
		return errors.New("one or more user ids not deleted")
	}
	return nil
}
