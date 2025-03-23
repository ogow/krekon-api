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
		a.handleGetHttp(w, r)
	case http.MethodPost:
		a.handlePostHttp(w, r)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

func (a *Api) handleGetHttp(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	result, err := a.db.GetHttpEntries(q)
	if err != nil {
		http.Error(w, "Could not get http entries", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Could not json encode http entries", http.StatusInternalServerError)
		return
	}
}

func (a *Api) handlePostHttp(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var result *db.HttpInfoContract
	if err := d.Decode(&result); err != nil {
		http.Error(w, "could not json decode req body", http.StatusBadRequest)
		return
	}

	_, err := a.db.StoreHttpEntries(result)
	if err != nil {
		http.Error(w, "Could not get entries", http.StatusInternalServerError)
		return
	}
}
