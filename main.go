package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sod-lol/url-shortener/core/repos/url_repo"
	"github.com/sod-lol/url-shortener/routers"
	"github.com/sod-lol/url-shortener/services/redis"
)

func main() {

	router := gin.New()

	root := context.Background()

	redisClient := redis.CreateRedisClient("redis-auth:6379", "", 0)
	defer redisClient.Close()
	urlRepo := url_repo.URLRedisStorage{
		Client: redisClient,
	}

	ctxWithRepo := context.WithValue(root, "user-repo", urlRepo)

	routers.InitRoutes(ctxWithRepo, &router.RouterGroup)

	router.Run(":8080")
}
