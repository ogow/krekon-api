package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ogow/krekon-api/db"
)

func (a *Api) HandleHttpEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHttpEntries(w, r, a.db)
	case http.MethodPost:
		postHttpEntry(w, r, a.db)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

func (a *Api) HandleHttpEntryByHostName(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHttpEntriesByHostName(w, r, a.db)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

func getHttpEntriesByHostName(w http.ResponseWriter, r *http.Request, db *db.Db) {
	host := r.PathValue("host")
	result, err := db.GetHttpEntries(host)
	if err != nil {
		http.Error(w, "Could not get http entries", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Could not json encode http entries", http.StatusInternalServerError)
		return
	}
}

func getHttpEntries(w http.ResponseWriter, r *http.Request, db *db.Db) {
	q := r.URL.Query().Get("q")
	result, err := db.GetHttpEntries(q)
	if err != nil {
		http.Error(w, "Could not get http entries", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Could not json encode http entries", http.StatusInternalServerError)
		return
	}
}

func postHttpEntry(w http.ResponseWriter, r *http.Request, dbm *db.Db) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var result *db.HttpInfoContract
	if err := d.Decode(&result); err != nil {
		http.Error(w, "could not json decode req body", http.StatusBadRequest)
		return
	}

	httpId, err := dbm.StoreHttpEntry(result)
	if err != nil {
		http.Error(w, "Could not get entries", http.StatusInternalServerError)
		return
	}
	_, err = dbm.StoreHttpRef(httpId, result.Host)
	if err != nil {
		http.Error(w, "could not store dns ref in entries collection", http.StatusInternalServerError)
	}
}
