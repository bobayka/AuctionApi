package main

import (
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"math/rand"
	"time"

	"gitlab.com/bobayka/courseproject/cmd/auth-api"
	"gitlab.com/bobayka/courseproject/cmd/lot-api"
	"gitlab.com/bobayka/courseproject/cmd/web-api"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"log"
	"net/http"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	db, err := postgres.PGInit("localhost", 5432, "bobayka", "12345", "TinkoffProj")
	if err != nil {
		log.Fatalf("pginit: %s", err)
	}
	StmtsStorage, err := postgres.NewUsersStorage(db)
	if err != nil {
		log.Fatalf("error in creation usersstorage: %s", err)
	}
	defer StmtsStorage.Close()

	r := chi.NewRouter()

	auth := authhandlers.NewAuthHandler(StmtsStorage)
	ra := auth.Routes()

	lotServ := lothandlers.NewLotServiceHandler(StmtsStorage)
	rl := lotServ.Routes()

	webServ := webapi.NewWebHandler(StmtsStorage)
	rw := webServ.Routes()

	r.Mount("/v1/auction", ra)
	r.Mount("/v1/auction/lots", rl)
	r.Mount("/w", rw)

	log.Fatal(http.ListenAndServe(":5000", r))
}
