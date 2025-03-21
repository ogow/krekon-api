package db

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OpenPorts struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Host      string        `bson:"host,omitempty" json:"host,omitempty"`
	IP        string        `bson:"ip,omitempty" json:"ip,omitempty"`
	Ports     []int         `bson:"ports,omitempty" json:"ports,omitempty"`
	Timestamp time.Time     `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}

func (db *Db) GetHostEntries(r string) ([]*OpenPorts, error) {
	collection := db.mongoClient.Database(db.name).Collection("hosts")

	filter := bson.D{{
		"host",
		bson.D{{
			"$regex", r,
		}},
	}}

	cur, err := collection.Find(db.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find any documents in hosts collection, err: %v", err)
	}

	var results []*OpenPorts
	if err := cur.All(db.ctx, &results); err != nil {
		return nil, fmt.Errorf("failed parse host documents, err: %v", err)
	}

	if len(results) == 0 {
		return []*OpenPorts{}, err
	}

	return results, nil
}
