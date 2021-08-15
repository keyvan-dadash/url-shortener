package redis

import (
	"github.com/go-redis/redis/v8"
)

//CreateRedisClient return Redis structure with given options
func CreateRedisClient(Addr string, Password string, DB int) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis-auth:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return redisClient
}
