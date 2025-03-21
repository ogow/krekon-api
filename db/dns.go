package db

import (
	"fmt"
	"time"

	"github.com/projectdiscovery/retryabledns"
	"go.mongodb.org/mongo-driver/v2/bson"
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

// get all dns entries based on a regex
func (db *Db) GetDnsEntries(r string) ([]*ShitBrokenDnsxPackage, error) {
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

	var results []*ShitBrokenDnsxPackage
	if err := cur.All(db.ctx, &results); err != nil {
		return nil, fmt.Errorf("failed parse dns documents, err: %v", err)
	}

	if len(results) == 0 {
		return []*ShitBrokenDnsxPackage{}, err
	}

	return results, nil
}
