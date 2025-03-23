package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ogow/krekon-api/db"
)

func (a *Api) HandleTlsEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.getTlsEntries(w, r)
	case http.MethodPost:
		a.postTlsEntry(w, r)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

func (a *Api) HandleTlsEntry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.getTlsEntry(w, r)
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
		return
	}
}

func (a *Api) getTlsEntry(w http.ResponseWriter, r *http.Request) {
	hostname := r.PathValue("host")

	result, err := a.db.GetTlsEntry(hostname)
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

func (a *Api) getTlsEntries(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	result, err := a.db.GetTlsEntries(q)
	if err != nil {
		http.Error(w, "Could not get tls entries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "Could not json encode entries", http.StatusInternalServerError)
		return
	}
}

func (a *Api) postTlsEntry(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var result *db.TlsContract
	if err := d.Decode(&result); err != nil {
		http.Error(w, "could not json decode req body", http.StatusBadRequest)
		return
	}

	tlsId, err := a.db.StoreTlsEntry(result)
	if err != nil {
		http.Error(w, "Could not get entries", http.StatusInternalServerError)
		return
	}

	_, err = a.db.StoreTlsRef(tlsId, result.Host)
	if err != nil {
		http.Error(w, "Could not store tls ref in entries collection", http.StatusInternalServerError)
		return
	}
}
