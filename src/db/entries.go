package db

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type EntryContract struct {
	Type      string          `bson:"-" json:"type,omitempty"`
	ID        bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	Host      string          `bson:"host,omitempty" json:"host,omitempty"`
	Dns       []bson.ObjectID `bson:"dns,omitempty" json:"dns,omitempty"`
	Tls       []bson.ObjectID `bson:"tls,omitempty" json:"tls,omitempty"`
	Http      []bson.ObjectID `bson:"http,omitempty" json:"http,omitempty"`
	Hosts     []bson.ObjectID `bson:"hosts,omitempty" json:"hosts,omitempty"`
	CreatedAt time.Time       `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

type EntryContractDetailed struct {
	Type      string              `bson:"-" json:"type,omitempty"`
	ID        bson.ObjectID       `bson:"_id,omitempty" json:"id"`
	Host      string              `bson:"host,omitempty" json:"host,omitempty"`
	Dns       []DnsContract       `bson:"dns,omitempty" json:"dns,omitempty"`
	Tls       []TlsContract       `bson:"tls,omitempty" json:"tls,omitempty"`
	Hosts     []OpenPortsContract `bson:"hosts,omitempty" json:"hosts,omitempty"`
	Http      []HttpInfoContract  `bson:"http,omitempty" json:"http,omitempty"`
	CreatedAt time.Time           `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

func (db *Db) DeleteEntry(hostname string) error {
	collEntries := db.mongoClient.Database(db.name).Collection("entries")
	collDns := db.mongoClient.Database(db.name).Collection("dns")
	collTls := db.mongoClient.Database(db.name).Collection("tls")
	collHosts := db.mongoClient.Database(db.name).Collection("hosts")
	collHttp := db.mongoClient.Database(db.name).Collection("http")

	// Starts a session on the client
	// Defers ending the session after the transaction is committed or ended

	var entry EntryContract
	if err := collEntries.FindOneAndDelete(db.ctx, bson.D{{"host", hostname}}).Decode(&entry); err != nil {
		return fmt.Errorf("could not find and delete document, err: %v", err)
	}

	// these delete statments should probalby be rewritten, return err at end of all delete statments
	if len(entry.Dns) > 0 {
		if _, err := collDns.DeleteMany(db.ctx, bson.M{"_id": bson.M{"$in": entry.Dns}}); err != nil {
			return err
		}
	}
	if len(entry.Tls) > 0 {
		if _, err := collTls.DeleteMany(db.ctx, bson.M{"_id": bson.M{"$in": entry.Tls}}); err != nil {
			return err
		}
	}
	if len(entry.Http) > 0 {
		if _, err := collHttp.DeleteMany(db.ctx, bson.M{"_id": bson.M{"$in": entry.Http}}); err != nil {
			return err
		}
	}
	if len(entry.Hosts) > 0 {
		if _, err := collHosts.DeleteMany(db.ctx, bson.M{"_id": bson.M{"$in": entry.Hosts}}); err != nil {
			return err
		}
	}

	return nil
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

func (db *Db) GetEntriesByHostNameDetailed(hostname string) ([]*EntryContractDetailed, error) {
	coll := db.mongoClient.Database(db.name).Collection("entries")

	pipeline := bson.A{
		bson.M{"$match": bson.M{"host": hostname}},
		bson.M{"$lookup": bson.M{
			"from":         "hosts",
			"localField":   "hosts",
			"foreignField": "_id",
			"as":           "hosts",
		}},
		bson.M{"$lookup": bson.M{
			"from":         "tls",
			"localField":   "tls",
			"foreignField": "_id",
			"as":           "tls",
		}},
		bson.M{"$lookup": bson.M{
			"from":         "http",
			"localField":   "http",
			"foreignField": "_id",
			"as":           "http",
		}},
		bson.M{"$lookup": bson.M{
			"from":         "dns",
			"localField":   "dns",
			"foreignField": "_id",
			"as":           "dns",
		}},
	}

	cursor, err := coll.Aggregate(db.ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("could not aggregate pipeline, err: %v", err)
	}

	var results []*EntryContractDetailed
	if err := cursor.All(db.ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to parse entries detailed result, err: %v", err)
	}

	if len(results) == 0 {
		return []*EntryContractDetailed{}, err
	}

	return results, nil
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

// get all entries based on host name
func (db *Db) GetEntriesDetailed(r string) ([]*EntryContractDetailed, error) {
	coll := db.mongoClient.Database(db.name).Collection("entries")

	pipeline := bson.A{
		bson.M{"$match": bson.M{"host": bson.M{"$regex": r}}},
		bson.M{"$lookup": bson.M{
			"from":         "hosts",
			"localField":   "hosts",
			"foreignField": "_id",
			"as":           "hosts",
		}},
		bson.M{"$lookup": bson.M{
			"from":         "tls",
			"localField":   "tls",
			"foreignField": "_id",
			"as":           "tls",
		}},
		bson.M{"$lookup": bson.M{
			"from":         "http",
			"localField":   "http",
			"foreignField": "_id",
			"as":           "http",
		}},
		bson.M{"$lookup": bson.M{
			"from":         "dns",
			"localField":   "dns",
			"foreignField": "_id",
			"as":           "dns",
		}},
	}

	cursor, err := coll.Aggregate(db.ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("could not aggregate pipeline, err: %v", err)
	}

	var results []*EntryContractDetailed
	if err := cursor.All(db.ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to parse entries detailed result, err: %v", err)
	}

	if len(results) == 0 {
		return []*EntryContractDetailed{}, err
	}

	return results, nil
}
