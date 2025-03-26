package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ogow/krekon-api/db"
)

func (a *Api) HandleHostEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHostEntries(w, r, a.db)
	case http.MethodPost:
		postHostEntry(w, r, a.db)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

func (a *Api) HandleHostEntry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHostEntry(w, r, a.db)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

func getHostEntries(w http.ResponseWriter, r *http.Request, db *db.Db) {
	q := r.URL.Query().Get("q")
	result, err := db.GetHostEntries(q)
	if err != nil {
		http.Error(w, "Could not get host entries", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Could not json encode host entries", http.StatusInternalServerError)
	}
}

func getHostEntry(w http.ResponseWriter, r *http.Request, db *db.Db) {
	host := r.PathValue("host")
	result, err := db.GetHostEntry(host)
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

func postHostEntry(w http.ResponseWriter, r *http.Request, dba *db.Db) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var result db.OpenPortsContract
	if err := d.Decode(&result); err != nil {
		http.Error(w, "could not json decode req body", http.StatusBadRequest)
		return
	}

	hostId, err := dba.StoreHostEntry(result)
	if err != nil {
		http.Error(w, "Could not get entries", http.StatusInternalServerError)
		return
	}

	_, err = dba.StoreHostsRef(hostId, result.Host)
	if err != nil {
		http.Error(w, "could not store hosts ref in entries collection", http.StatusInternalServerError)
	}
}
