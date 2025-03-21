package db

import (
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// get all dns entries based on a regex
type HttpInfoContract struct {
	ID        bson.ObjectID       `bson:"_id,omitempty"`
	Url       string              `bson:"url,omitempty"`
	Code      int                 `bson:"code,omitempty"`
	Title     string              `bson:"title,omitempty"`
	Length    int64               `bson:"length,omitempty"`
	Words     int                 `bson:"words,omitempty"`
	Location  string              `bson:"location,omitempty"`
	Port      string              `bson:"port,omitempty"`
	Headers   http.Header         `bson:"headers,omitempty"`
	Body      string              `bson:"body,omitempty"`
	Tech      map[string]struct{} `bson:"tech,omitempty"`
	Timestamp time.Time           `json:"timestamp,omitempty"`
}

func (db *Db) GetHttpEntries(r string) ([]*HttpInfoContract, error) {
	collection := db.mongoClient.Database(db.name).Collection("http")

	filter := bson.D{{
		"url",
		bson.D{{
			"$regex", r,
		}},
	}}

	cur, err := collection.Find(db.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find any documents in http collection, err: %v", err)
	}

	var results []*HttpInfoContract
	if err := cur.All(db.ctx, &results); err != nil {
		return nil, fmt.Errorf("failed parse dns documents, err: %v", err)
	}

	if len(results) == 0 {
		return []*HttpInfoContract{}, err
	}

	return results, nil
}
