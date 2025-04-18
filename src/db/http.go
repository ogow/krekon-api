package db

import (
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// get all dns entries based on a regex
type HttpInfoContract struct {
	ID        bson.ObjectID       `bson:"_id,omitempty" json:"id,omitempty"`
	Host      string              `bson:"host,omitempty" json:"host,omitempty"`
	Url       string              `bson:"url,omitempty" json:"url,omitempty"`
	Code      int                 `bson:"code,omitempty" json:"code,omitempty"`
	Title     string              `bson:"title,omitempty" json:"title,omitempty"`
	Length    int64               `bson:"length,omitempty" json:"length,omitempty"`
	Words     int                 `bson:"words,omitempty" json:"words,omitempty"`
	Location  string              `bson:"location,omitempty" json:"location,omitempty"`
	Port      string              `bson:"port,omitempty" json:"port,omitempty"`
	Headers   http.Header         `bson:"headers,omitempty" json:"headers,omitempty"`
	Body      string              `bson:"body,omitempty" json:"body,omitempty"`
	Tech      map[string]struct{} `bson:"tech,omitempty" json:"tech,omitempty"`
	Timestamp time.Time           `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
}

func (db *Db) StoreHttpRef(httpRefId interface{}, hostname string) (interface{}, error) {
	entriesCollection := db.mongoClient.Database(db.name).Collection("entries")
	update := bson.M{
		"$addToSet": bson.M{"http": httpRefId},
		"$set":      bson.M{"created_at": time.Now()},
	}
	opts := options.UpdateOne().SetUpsert(true)

	filter := bson.M{"host": hostname}
	// update := bson.M{"$set": bson.M{}}

	result, err := entriesCollection.UpdateOne(db.ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return result, nil
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

func (db *Db) GetHttpEntriesByHostName(host string) ([]HttpInfoContract, error) {
	collection := db.mongoClient.Database(db.name).Collection("http")
	m := fmt.Sprintf(`^%s$`, host)
	filter := bson.D{{
		"host", bson.D{{"$regex", m}},
	}}

	cur, err := collection.Find(db.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find any documents in http collection, err: %v", err)
	}

	var results []HttpInfoContract
	if err := cur.All(db.ctx, &results); err != nil {
		return nil, fmt.Errorf("failed parse http documents, err: %v", err)
	}

	if len(results) == 0 {
		return []HttpInfoContract{}, err
	}

	return results, nil
}

func (db *Db) StoreHttpEntry(httpInfo *HttpInfoContract) (interface{}, error) {
	collection := db.mongoClient.Database(db.name).Collection("http")
	// get subdomain name from url
	filter := bson.M{"url": httpInfo.Url}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	upd := bson.M{
		"$set": httpInfo,
		// "$set": bson.M{fmt.Sprintf("http.%s", p.Port): p},
		//"$addToSet": bson.M{"ports": portInt},
	}

	var result *HttpInfoContract
	if err := collection.FindOneAndUpdate(db.ctx, filter, upd, opts).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to find and update http document for %v, err:%v", httpInfo, err)
	}

	return result.ID, nil
}
