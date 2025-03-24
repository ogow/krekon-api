package db

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type EntryContract struct {
	ID        bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	Host      string          `bson:"host,omitempty" json:"host,omitempty"`
	Dns       []bson.ObjectID `bson:"dns,omitempty" json:"dns,omitempty"`
	Tls       []bson.ObjectID `bson:"tls,omitempty" json:"tls,omitempty"`
	Hosts     []bson.ObjectID `bson:"hosts,omitempty" json:"hosts,omitempty"`
	CreatedAt time.Time       `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

func (db *Db) StoreEntry(entry EntryContract) (interface{}, error) {
	collection := db.mongoClient.Database(db.name).Collection("entries")

	opts := options.Find().SetSort(bson.D{{"created_at", -1}}) // sort decending date, latest date first
	// check id host already in db
	cur, err := collection.Find(db.ctx, bson.M{"host": entry.Host}, opts)
	if err != nil {
		return nil, err
	}

	// maybe this part should be rebuilt, i could just use findOne im guessing
	var results []EntryContract
	if err := cur.All(db.ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		entry := EntryContract{
			Host:      entry.Host,
			CreatedAt: time.Now(),
		}

		ins, err := collection.InsertOne(db.ctx, entry)
		if err != nil {
			return nil, fmt.Errorf("failed to store host: %s, err: %v", entry.Host, err)
		}
		return ins.InsertedID, nil
	}
	if len(results) > 1 {
		return nil, fmt.Errorf("host %v has duplicate entries in mongodb collection entries", entry.Host)
	}
	return results[0].ID, nil
}

func (db *Db) GetEntry(hostname string) (*EntryContract, error) {
	coll := db.mongoClient.Database(db.name).Collection("entries")
	filter := bson.D{{"host", hostname}}
	// Retrieves the first matching document
	var result *EntryContract
	if err := coll.FindOne(db.ctx, filter).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to find one document for %s, err: %v", hostname, err)
	}
	return result, nil
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
