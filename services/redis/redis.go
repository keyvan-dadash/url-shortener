package redis

import (
	"github.com/go-redis/redis/v8"
)

//CreateRedisClient return Redis structure with given options
func CreateRedisClient(Addr string, Password string, DB int) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: Password, // no password set
		DB:       DB,       // use default DB
	})

	return redisClient
}
