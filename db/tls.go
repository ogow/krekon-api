package db

import (
	"fmt"
	"time"

	"github.com/projectdiscovery/tlsx/pkg/tlsx/clients"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TlsContract struct {
	// Timestamp is the timestamp for certificate response
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Timestamp *time.Time    `bson:"timestamp,omitempty"`
	// Host is the host to make request to
	Host string `bson:"host"`
	// IP is the IP address the request was made to
	IP string `bson:"ip,omitempty"`
	// Port is the port to make request to
	Port string `bson:"port"`
	// ProbeStatus is false if the tls probe failed
	ProbeStatus bool `bson:"probe_status"`
	// Version is the tls version responded by the server
	Version string `bson:"tls_version,omitempty"`
	// Cipher is the cipher for the tls request
	Cipher string `bson:"cipher,omitempty"`
	// CertificateResponse is the leaf certificate embedded in json
	*clients.CertificateResponse `bson:",inline"`
	// TLSConnection is the client used for TLS connection
	// when ran using scan-mode auto.
	TLSConnection string `bson:"tls_connection,omitempty"`
	// Chain is the chain of certificates
	Chain              []*clients.CertificateResponse `bson:"chain,omitempty"`
	JarmHash           string                         `bson:"jarm_hash,omitempty"`
	Ja3Hash            string                         `bson:"ja3_hash,omitempty"`
	Ja3sHash           string                         `bson:"ja3s_hash,omitempty"`
	ServerName         string                         `bson:"sni,omitempty"`
	VersionEnum        []string                       `bson:"version_enum,omitempty"`
	TlsCiphers         []clients.TlsCiphers           `bson:"cipher_enum,omitempty"`
	ClientCertRequired *bool                          `bson:"client_cert_required,omitempty"`
	// ClientHello        *ztls.ClientHello              `json:"client_hello,omitempty"`
	// ServerHello        *ztls.ServerHello              `json:"servers_hello,omitempty"`
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
