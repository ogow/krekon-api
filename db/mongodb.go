package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Connect to MongoDB
func ConnectMongo(ctx context.Context, conn string, uname, passw string) *mongo.Client {
	credentials := options.Credential{
		Username: uname,
		Password: passw,
	}
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credentials)
	clientOptions := options.Client().ApplyURI(conn).SetAuth(credentials)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	ctxConn, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	err = client.Ping(ctxConn, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
