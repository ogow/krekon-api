package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "this is the root\n")
}

func (a *Api) GetEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		q := r.URL.Query().Get("q")
		result, err := a.db.GetEntries(q)
		if err != nil {
			http.Error(w, "Could not get entries", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Could not json encode entries", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
	}
}

func (a *Api) GetDnsEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		q := r.URL.Query().Get("q")
		result, err := a.db.GetDnsEntries(q)
		if err != nil {
			http.Error(w, "Could not get entries", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Could not json encode entries", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
	}
}

func (a *Api) GetTlsEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		q := r.URL.Query().Get("q")
		result, err := a.db.GetTlsEntries(q)
		if err != nil {
			http.Error(w, "Could not get tls entries", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Could not json encode entries", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
	}
}

func (a *Api) GetHttpEntries(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		q := r.URL.Query().Get("q")
		result, err := a.db.GetHttpEntries(q)
		if err != nil {
			http.Error(w, "Could not get http entries", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Could not json encode http entries", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprint("http method not supported"), http.StatusBadRequest)
	}
}
