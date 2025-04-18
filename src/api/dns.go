package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ogow/krekon-api/db"
)

func (a *Api) HandleDnsEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.handleGetDns(w, r)
	case http.MethodPost:
		a.handlePostDns(w, r)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

func (a *Api) HandleDnsEntry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.handleGetDnsSingleHost(w, r)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

// func (a *Api) HandleDnsEntryId(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		// a.handleGetDnsSingleHost(w, r)
// 	case http.MethodPost:
// 	default:
// 		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
// 		return
// 	}
// }

func (a *Api) handleGetDns(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	result, err := a.db.GetDnsEntries(q)
	if err != nil {
		http.Error(w, "Could not get entries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Could not json encode entries", http.StatusInternalServerError)
		return
	}
}

func (a *Api) handleGetDnsSingleHost(w http.ResponseWriter, r *http.Request) {
	host := r.PathValue("host")
	result, err := a.db.GetDnsEntriesByHostName(host)
	if err != nil {
		http.Error(w, "Could not get entries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Could not json encode entries", http.StatusInternalServerError)
		return
	}
}

func (a *Api) handlePostDns(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var result db.ShitBrokenDnsxPackage
	if err := d.Decode(&result); err != nil {
		http.Error(w, "could not json decode req body", http.StatusBadRequest)
		return
	}

	dnsId, err := a.db.StoreDnsEntry(result)
	if err != nil {
		if strings.Contains(err.Error(), "dns record not changed") {
			fmt.Fprintf(w, "dns record already exists and has not changed\n")
			return
		}
		http.Error(w, "Could not get entries", http.StatusInternalServerError)
		return
	}

	_, err = a.db.StoreDnsRef(dnsId, result.Host)
	if err != nil {
		http.Error(w, "could not store dns ref in entries collection", http.StatusInternalServerError)
	}
	// also store a ref in entries collection
}
