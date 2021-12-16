package main

import (
	"log"

	"github.com/redis-developer/basic-redis-leaderboard-demo-go/api"
	"github.com/redis-developer/basic-redis-leaderboard-demo-go/config"
	"github.com/redis-developer/basic-redis-leaderboard-demo-go/controller"
	"github.com/redis-developer/basic-redis-leaderboard-demo-go/internal"
	"github.com/redis-developer/basic-redis-leaderboard-demo-go/redis"
)

func main() {
	newConfig := config.NewConfig()

	options, err := redis.NewOptions(newConfig.Redis)
	if err != nil {
		log.Fatalln(err)
	}

	newRedis := redis.New(options)

	newServer := api.NewServer(newConfig.Api)

	controller.SetRedis(newRedis)

	err = controller.ImportCompanies(newConfig.Import.Path(), newRedis)
	if err != nil {
		log.Fatalln(err)
	}

	internal.Waiting(newServer, newRedis)

	log.Println("Server exiting")

}
