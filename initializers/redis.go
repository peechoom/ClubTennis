package initializers

import (
	"ClubTennis/config"
	"fmt"

	"github.com/go-redis/redis"
)

func GetClient(c *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Pass,
		DB:       1,
	})
	return client
}

func GetTestClient(c *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Pass,
		DB:       14,
	})
	return client
}
