package redis

import (
	"github.com/go-redis/redis"
	"github.com/pauljamescleary/gomin/pkg/common/config"
)

func StartRedisClient(c config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.RedisUrl,
		Password: "",
		DB:       0,
	})
	return client, nil
}
