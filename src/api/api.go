package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ogow/krekon-api/db"
)

type Api struct {
	ctx context.Context
	db  *db.Db
}

type ApiOpts struct {
	Ctx context.Context
	Db  *db.Db
}

func New(opts ApiOpts) *Api {
	return &Api{
		db:  opts.Db,
		ctx: opts.Ctx,
	}
}

func (a *Api) ServeApi() {
	http.HandleFunc("/{$}", HandleRoot)
	http.HandleFunc("/entries", a.HandleEntries)
	http.HandleFunc("/entries/{host}", a.HandleEntry)
	http.HandleFunc("/hosts", a.HandleHostEntries)
	http.HandleFunc("/hosts/{host}", a.HandleHostEntry)
	http.HandleFunc("/dns", a.HandleDnsEntries)
	http.HandleFunc("/dns/{host}", a.HandleDnsEntry)
	http.HandleFunc("/dns/{host}/{id}", a.HandleDnsEntry)
	http.HandleFunc("/tls", a.HandleTlsEntries)
	http.HandleFunc("/tls/{host}", a.HandleTlsEntry)
	http.HandleFunc("/http", a.HandleHttpEntries)
	http.HandleFunc("/http/{host}", a.HandleHttpEntryByHostName)

	fmt.Println("serving API on port", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
