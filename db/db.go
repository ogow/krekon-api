package db

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Db struct {
	ctx         context.Context
	mongoClient *mongo.Client
	redisClient *redis.Client
	name        string
}

// NewHTTPService creates a new HTTPService
func New(mongoClient *mongo.Client, redisClient *redis.Client, name string) *Db {
	return &Db{
		name:        name,
		mongoClient: mongoClient,
		redisClient: redisClient,
	}
}

func (db *Db) Stop() {
	db.redisClient.Close()
	db.mongoClient.Disconnect(db.ctx)
}
