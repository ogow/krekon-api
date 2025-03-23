package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ogow/krekon-api/db"
)

func (a *Api) HandleEntry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.handleGetEntry(w, r)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

func (a *Api) HandleEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.handleGetEntries(w, r)
	case http.MethodPost:
		a.handlePostEntry(w, r)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

func (a *Api) handleGetEntries(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	result, err := a.db.GetEntries(q)
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

func (a *Api) handlePostEntry(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	var result db.EntryContract
	if err := d.Decode(&result); err != nil {
		http.Error(w, "could not json decode req body", http.StatusBadRequest)
		return
	}

	_, err := a.db.StoreEntry(result)
	if err != nil {
		http.Error(w, "Could not get entries", http.StatusInternalServerError)
		return
	}
}

func (a *Api) handleGetEntry(w http.ResponseWriter, r *http.Request) {
	hostname := r.PathValue("host")

	result, err := a.db.GetEntry(hostname)
	if err != nil {
		http.Error(w, "Could not get entry", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Could not json encode entry", http.StatusInternalServerError)
		return
	}
}
