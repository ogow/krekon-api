package db

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EntryContract struct {
	ID        bson.ObjectID   `bson:"_id,omitempty"`
	Host      string          `bson:"host,omitempty"`
	Dns       []bson.ObjectID `bson:"dns,omitempty"`
	Tls       []bson.ObjectID `bson:"tls,omitempty"`
	Hosts     []bson.ObjectID `bson:"hosts,omitempty"`
	CreatedAt time.Time       `bson:"created_at,omitempty"`
}

// get all entries based on host name
func (db *Db) GetEntries(r string) ([]*EntryContract, error) {
	collection := db.mongoClient.Database(db.name).Collection("entries")

	filter := bson.D{{
		"host",
		bson.D{{
			"$regex", r,
		}},
	}}

	cur, err := collection.Find(db.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find any documents in entries collection, err: %v", err)
	}

	var results []*EntryContract
	if err := cur.All(db.ctx, &results); err != nil {
		return nil, fmt.Errorf("failed parse entries documents, err: %v", err)
	}

	if len(results) == 0 {
		return []*EntryContract{}, err
	}

	return results, nil
}
