package redis

import (
	"github.com/redis/go-redis/v9"
)

type Client = redis.Client

func RedisInit(redisUrl string) *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    redisUrl,
	})
	return redis
}
