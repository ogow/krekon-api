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
		if r.URL.Query().Get("detailed") == "true" {
			a.getEntriesByHostNameDetailed(w, r)
		} else {
			a.getEntry(w, r)
		}
	case http.MethodDelete:
		a.deleteEntry(w, r)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

func (a *Api) HandleEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Query().Get("detailed") == "true" {
			a.getEntriesDetailed(w, r)
		} else {
			a.getEntries(w, r)
		}

	case http.MethodPost:
		a.postEntry(w, r)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

func (a *Api) deleteEntry(w http.ResponseWriter, r *http.Request) {
	hostname := r.PathValue("host")
	if err := a.db.DeleteEntry(hostname); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete entry %s, err %v", hostname, err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (a *Api) getEntries(w http.ResponseWriter, r *http.Request) {
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

func (a *Api) getEntriesDetailed(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	result, err := a.db.GetEntriesDetailed(q)
	if err != nil {
		e := fmt.Sprintf("Could not get detailed entries, err: %v", err)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Could not json encode entries", http.StatusInternalServerError)
		return
	}
}

func (a *Api) getEntriesByHostNameDetailed(w http.ResponseWriter, r *http.Request) {
	hostname := r.PathValue("host")

	result, err := a.db.GetEntriesByHostNameDetailed(hostname)
	if err != nil {
		e := fmt.Sprintf("Could not get detailed entries, err: %v", err)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Could not json encode entries", http.StatusInternalServerError)
		return
	}
}

func (a *Api) postEntry(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusCreated)
}

func (a *Api) getEntry(w http.ResponseWriter, r *http.Request) {
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
