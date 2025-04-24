package db

import (
	"fmt"
	"time"

	"github.com/projectdiscovery/tlsx/pkg/tlsx/clients"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TlsContract struct {
	// Timestamp is the timestamp for certificate response
	Type      string        `bson:"-" json:"type,omitempty"`
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp *time.Time    `bson:"timestamp,omitempty" json:"timestamp,omitempty"`
	// Host is the host to make request to
	Host string `bson:"host" json:"host"`
	// IP is the IP address the request was made to
	IP string `bson:"ip,omitempty" json:"ip,omitempty"`
	// Port is the port to make request to
	Port string `bson:"port" json:"port"`
	// ProbeStatus is false if the TLS probe failed
	ProbeStatus bool `bson:"probe_status" json:"probe_status"`
	// Version is the TLS version responded by the server
	Version string `bson:"tls_version,omitempty" json:"tls_version,omitempty"`
	// Cipher is the cipher for the TLS request
	Cipher string `bson:"cipher,omitempty" json:"cipher,omitempty"`
	// CertificateResponse is the leaf certificate embedded in JSON
	*clients.CertificateResponse `bson:",inline" json:",inline"`
	// TLSConnection is the client used for TLS connection when ran using scan-mode auto
	TLSConnection string `bson:"tls_connection,omitempty" json:"tls_connection,omitempty"`
	// Chain is the chain of certificates
	Chain              []*clients.CertificateResponse `bson:"chain,omitempty" json:"chain,omitempty"`
	JarmHash           string                         `bson:"jarm_hash,omitempty" json:"jarm_hash,omitempty"`
	Ja3Hash            string                         `bson:"ja3_hash,omitempty" json:"ja3_hash,omitempty"`
	Ja3sHash           string                         `bson:"ja3s_hash,omitempty" json:"ja3s_hash,omitempty"`
	ServerName         string                         `bson:"sni,omitempty" json:"sni,omitempty"`
	VersionEnum        []string                       `bson:"version_enum,omitempty" json:"version_enum,omitempty"`
	TlsCiphers         []clients.TlsCiphers           `bson:"cipher_enum,omitempty" json:"cipher_enum,omitempty"`
	ClientCertRequired *bool                          `bson:"client_cert_required,omitempty" json:"client_cert_required,omitempty"`
}

func (db *Db) StoreTlsRef(tlsId interface{}, hostname string) (interface{}, error) {
	entriesCollection := db.mongoClient.Database(db.name).Collection("entries")
	// update := bson.D{{bson.D{{"$push", "$set", bson.D{{"dns", dnsRefId}}}}}}
	update := bson.M{
		"$addToSet": bson.M{"tls": tlsId},
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

func (db *Db) GetTlsEntry(hostname string) ([]*TlsContract, error) {
	collection := db.mongoClient.Database(db.name).Collection("tls")

	filter := bson.D{{
		"host", hostname,
	}}

	cur, err := collection.Find(db.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find any documents for %s in tls collection, err: %v", hostname, err)
	}

	var results []*TlsContract
	if err := cur.All(db.ctx, &results); err != nil {
		return nil, fmt.Errorf("failed parse tls documents, err: %v", err)
	}
	if len(results) == 0 {
		return []*TlsContract{}, err
	}

	return results, nil
}

func (db *Db) GetTlsEntries(r string) ([]*TlsContract, error) {
	collection := db.mongoClient.Database(db.name).Collection("tls")

	filter := bson.D{{
		"host",
		bson.D{{
			"$regex", r,
		}},
	}}

	cur, err := collection.Find(db.ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find any documents in tls collection, err: %v", err)
	}

	var results []*TlsContract
	if err := cur.All(db.ctx, &results); err != nil {
		return nil, fmt.Errorf("failed parse tls documents, err: %v", err)
	}
	if len(results) == 0 {
		return []*TlsContract{}, err
	}

	return results, nil
}

func (db *Db) StoreTlsEntry(tlsData TlsContract) (interface{}, error) {
	tlsCollection := db.mongoClient.Database(db.name).Collection("tls")

	cont := TlsContract{
		Timestamp:           tlsData.Timestamp,
		Host:                tlsData.Host,
		IP:                  tlsData.IP,
		Port:                tlsData.Port,
		ProbeStatus:         tlsData.ProbeStatus,
		Version:             tlsData.Version,
		Cipher:              tlsData.Cipher,
		CertificateResponse: tlsData.CertificateResponse,
		TLSConnection:       tlsData.TLSConnection,
		Chain:               tlsData.Chain,
		JarmHash:            tlsData.JarmHash,
		Ja3Hash:             tlsData.Ja3Hash,
		Ja3sHash:            tlsData.Ja3sHash,
		ServerName:          tlsData.ServerName,
		VersionEnum:         tlsData.VersionEnum,
		TlsCiphers:          tlsData.TlsCiphers,
		ClientCertRequired:  tlsData.ClientCertRequired,
	}

	// tlsFilter := bson.D{{"$set", bson.D{{ }}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	upd := bson.D{{"$set", cont}}
	tlsFilter := bson.M{"host": tlsData.Host}

	var tlsCon TlsContract

	err := tlsCollection.FindOneAndUpdate(db.ctx, tlsFilter, upd, opts).Decode(&tlsCon)
	if err != nil {
		return nil, err
	}
	return tlsCon.ID, nil
}
