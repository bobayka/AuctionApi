package main

import (
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

	auth := NewAuthHandler(StmtsStorage)
	r := auth.Routes()
	r.Mount("/v1/auction", r)
	log.Fatal(http.ListenAndServe(":5000", r))
}
