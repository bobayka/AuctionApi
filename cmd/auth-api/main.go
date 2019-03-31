package main

import (
	"github.com/go-chi/chi"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"log"
	"net/http"
)

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
	auth := NewAuthHandler(StmtsStorage)
	r = auth.Routes(r)
	lotServ := NewLotServiceHandler(StmtsStorage)
	r = lotServ.Routes(r)
	r.Mount("/v1/auction/", r)
	log.Fatal(http.ListenAndServe(":5000", r))
}
