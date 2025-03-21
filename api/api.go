package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
	http.HandleFunc("/entries", a.GetEntries)
	http.HandleFunc("/dns", a.GetDnsEntries)
	http.HandleFunc("/tls", a.GetTlsEntries)
	http.HandleFunc("/http", a.GetHttpEntries)

	fmt.Println("serving API on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
