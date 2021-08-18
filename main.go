package main

import (
	"log"
	"os"

	goRedis "github.com/go-redis/redis"

	"github.com/redis-developer/basic-redis-leaderboard-demo-go/api"
	"github.com/redis-developer/basic-redis-leaderboard-demo-go/config"
	"github.com/redis-developer/basic-redis-leaderboard-demo-go/controller"
	"github.com/redis-developer/basic-redis-leaderboard-demo-go/internal"
	"github.com/redis-developer/basic-redis-leaderboard-demo-go/redis"
)

const envRedisURL = "REDIS_URL"

func main() {
	newConfig := config.NewConfig()

	// ref https://pkg.go.dev/github.com/go-redis/redis?utm_source=gopls#ParseURL
	var newRedis *redis.Redis
	url, ok := os.LookupEnv(envRedisURL)
	if ok {
		opt, err := goRedis.ParseURL(url)
		if err != nil {
			log.Fatalln(err)
		}
		newRedis = redis.NewRedisFromOptions(opt)
	} else {
		newRedis = redis.NewRedis(newConfig.Redis)
	}

	newServer := api.NewServer(newConfig.Api)

	controller.SetRedis(newRedis)

	err := controller.ImportCompanies(newConfig.Import.Path(), newRedis)
	if err != nil {
		log.Fatalln(err)
	}

	internal.Waiting(newServer, newRedis)

	log.Println("Server exiting")

}
