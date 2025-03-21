package db

import (
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/projectdiscovery/retryabledns"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ShitBrokenDnsxPackage struct {
	Host           string                  `json:"host,omitempty"`
	TTL            uint32                  `json:"ttl,omitempty"`
	Resolver       []string                `json:"resolver,omitempty"`
	A              []string                `json:"a,omitempty"`
	AAAA           []string                `json:"aaaa,omitempty"`
	CNAME          []string                `json:"cname,omitempty"`
	MX             []string                `json:"mx,omitempty"`
	PTR            []string                `json:"ptr,omitempty"`
	SOA            []retryabledns.SOA      `json:"soa,omitempty"`
	NS             []string                `json:"ns,omitempty"`
	TXT            []string                `json:"txt,omitempty"`
	SRV            []string                `json:"srv,omitempty"`
	CAA            []string                `json:"caa,omitempty"`
	AllRecords     []string                `json:"all,omitempty"`
	Raw            string                  `json:"raw,omitempty"`
	HasInternalIPs bool                    `json:"has_internal_ips,omitempty"`
	InternalIPs    []string                `json:"internal_ips,omitempty"`
	StatusCode     string                  `json:"status_code,omitempty"`
	StatusCodeRaw  int                     `json:"status_code_raw,omitempty"`
	TraceData      *retryabledns.TraceData `json:"trace,omitempty"`
	AXFRData       *retryabledns.AXFRData  `json:"axfr,omitempty"`
	Timestamp      time.Time               `json:"timestamp,omitempty"`
	HostsFile      bool                    `json:"hosts_file,omitempty"`
}

type DnsContract struct {
	ID             bson.ObjectID          `bson:"_id,omitempty" json:"id"`
	Host           string                 `bson:"host,omitempty" json:"host,omitempty"`
	A              []string               `bson:"a,omitempty" json:"a,omitempty"`
	AAAA           []string               `bson:"aaaa,omitempty" json:"aaaa,omitempty"`
	CNAME          []string               `bson:"cname,omitempty" json:"cname,omitempty"`
	TXT            []string               `bson:"txt,omitempty" json:"txt,omitempty"`
	MX             []string               `bson:"mx,omitempty" json:"mx,omitempty"`
	NS             []string               `bson:"ns,omitempty" json:"ns,omitempty"`
	PTR            []string               `bson:"ptr,omitempty" json:"ptr,omitempty"`
	SOA            []retryabledns.SOA     `bson:"soa,omitempty" json:"soa,omitempty"`
	AXFRData       *retryabledns.AXFRData `bson:"axfrdata,omitempty" json:"axfrdata,omitempty"`
	StatusCode     string                 `bson:"statuscode,omitempty" json:"statuscode,omitempty"`
	HasInternalIPs bool                   `bson:"hasinternalips,omitempty" json:"hasinternalips,omitempty"`
	InternalIPs    []string               `bson:"internalips,omitempty" json:"internalips,omitempty"`
	Timestamp      time.Time              `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	AllRecords     []string               `bson:"allrecords,omitempty" json:"allrecords,omitempty"`
}

// get all dns entries based on a regex
func (db *Db) GetDnsEntries(r string) ([]*DnsContract, error) {
	collection := db.mongoClient.Database(db.name).Collection("dns")
	filter := bson.D{{
		"host",
		bson.D{{
			"$regex", r,
		}},
	}}

	cur, err := collection.Find(db.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find any documents in dns collection, err: %v", err)
	}

	var results []*DnsContract
	if err := cur.All(db.ctx, &results); err != nil {
		return nil, fmt.Errorf("failed parse dns documents, err: %v", err)
	}

	if len(results) == 0 {
		return []*DnsContract{}, err
	}

	return results, nil
}

func (db *Db) StoreDnsEntry(dnsData *ShitBrokenDnsxPackage) (interface{}, error) {
	dnsCollection := db.mongoClient.Database(db.name).Collection("dns")

	opts := options.Find().SetSort(bson.D{{"timestamp", -1}}) // sort decending date, latest date first

	findDnsRecords, err := dnsCollection.Find(db.ctx, bson.M{"host": dnsData.Host}, opts)
	if err != nil {
		return nil, err
	}

	var results []DnsContract
	if err := findDnsRecords.All(db.ctx, &results); err != nil {
		return nil, err
	}

	// insert dns record
	newDnsData := DnsContract{
		Host:           dnsData.Host,
		A:              dnsData.A,
		AAAA:           dnsData.AAAA,
		CNAME:          dnsData.CNAME,
		TXT:            dnsData.TXT,
		MX:             dnsData.MX,
		NS:             dnsData.NS,
		PTR:            dnsData.PTR,
		SOA:            dnsData.SOA,
		AXFRData:       dnsData.AXFRData,
		StatusCode:     dnsData.StatusCode,
		HasInternalIPs: dnsData.HasInternalIPs,
		InternalIPs:    dnsData.InternalIPs,
		Timestamp:      dnsData.Timestamp,
		AllRecords:     dnsData.AllRecords,
	}

	slices.Sort(newDnsData.A)
	slices.Sort(newDnsData.AAAA)
	slices.Sort(newDnsData.CNAME)
	slices.Sort(newDnsData.PTR)
	slices.Sort(newDnsData.NS)

	if len(results) > 0 {
		fdns := results[0]

		var dnsHasChanged bool = false

		switch {
		case !reflect.DeepEqual(newDnsData.AAAA, fdns.AAAA):
			dnsHasChanged = true
		case !reflect.DeepEqual(newDnsData.A, fdns.A):
			dnsHasChanged = true
		case !reflect.DeepEqual(newDnsData.CNAME, fdns.CNAME):
			dnsHasChanged = true
		case !reflect.DeepEqual(newDnsData.PTR, fdns.PTR):
			dnsHasChanged = true
		case !reflect.DeepEqual(newDnsData.NS, fdns.NS):
			dnsHasChanged = true
		case newDnsData.StatusCode != fdns.StatusCode:
			dnsHasChanged = true
		}
		if dnsHasChanged {
			id, err := dnsCollection.InsertOne(db.ctx, newDnsData)
			if err != nil {
				return nil, err
			}
			return id.InsertedID, nil
		}
	} else {
		id, err := dnsCollection.InsertOne(db.ctx, newDnsData)
		if err != nil {
			return nil, err
		}
		return id.InsertedID, nil
	}

	return nil, fmt.Errorf("dns record not changed for %v", newDnsData.Host)
}
