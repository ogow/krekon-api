package db

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type OpenPortsContract struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Host      string        `bson:"host,omitempty" json:"host,omitempty"`
	IP        string        `bson:"ip,omitempty" json:"ip,omitempty"`
	Ports     []int         `bson:"ports,omitempty" json:"ports,omitempty"`
	Timestamp time.Time     `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}

func (db *Db) StoreHostsRef(hostId interface{}, hostname string) (interface{}, error) {
	entriesCollection := db.mongoClient.Database(db.name).Collection("entries")
	// update := bson.D{{bson.D{{"$push", "$set", bson.D{{"dns", dnsRefId}}}}}}
	update := bson.M{
		"$addToSet": bson.M{"hosts": hostId},
		"$set":      bson.M{"created_at": time.Now()},
	}

	opts := options.UpdateOne().SetUpsert(true)

	filter := bson.M{"host": hostname}
	// update := bson.M{"$set": bson.M{}}

	id, err := entriesCollection.UpdateOne(db.ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return id, err
}

func (db *Db) StoreHostEntry(op OpenPortsContract) (interface{}, error) {
	collection := db.mongoClient.Database(db.name).Collection("hosts")

	// filter := bson.M{"$or": []bson.M{{"name": host}, {"ip": host}}}
	filter := bson.M{"ip": op.IP}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	// upd := bson.M{
	// 	// "$set": bson.M{"ports": result.Ports},
	// 	"$addToSet": bson.M{"ports": e},
	// }
	upd := bson.M{"$set": op}
	var ports OpenPortsContract
	if err := collection.FindOneAndUpdate(db.ctx, filter, upd, opts).Decode(&ports); err != nil {
		return nil, fmt.Errorf("failed to updated and find port document for %s, err: %v", op.Host, err)
	}

	return ports.ID, nil
}

func (db *Db) GetHostEntries(r string) ([]*OpenPortsContract, error) {
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

	var results []*OpenPortsContract
	if err := cur.All(db.ctx, &results); err != nil {
		return nil, fmt.Errorf("failed parse host documents, err: %v", err)
	}

	if len(results) == 0 {
		return []*OpenPortsContract{}, err
	}

	return results, nil
}

func (db *Db) GetHostEntry(hostname string) ([]*OpenPortsContract, error) {
	collection := db.mongoClient.Database(db.name).Collection("hosts")

	filter := bson.D{{
		"host", hostname,
	}}

	cur, err := collection.Find(db.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find any documents in hosts collection, err: %v", err)
	}

	var results []*OpenPortsContract
	if err := cur.All(db.ctx, &results); err != nil {
		return nil, fmt.Errorf("failed parse host documents, err: %v", err)
	}

	if len(results) == 0 {
		return []*OpenPortsContract{}, err
	}

	return results, nil
}
