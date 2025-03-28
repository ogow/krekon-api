package main

import (
	"context"

	"github.com/ogow/krekon-api/api"
	"github.com/ogow/krekon-api/db"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mongoClient := db.ConnectMongo(ctx, "mongodb://mongodb:27017", "root", "krekonApiPasswordOskarWasHere123")
	redisClient := db.ConnectRedis(ctx, "redis:6379")

	db := db.New(mongoClient, redisClient, "recondb")
	defer db.Stop()

	apiOpts := api.ApiOpts{
		Db:  db,
		Ctx: ctx,
	}

	api := api.New(apiOpts)

	api.ServeApi()
}
