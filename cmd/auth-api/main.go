package main

import (
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"log"
	"net/http"
)

func main() {
	_, err := postgres.PGInit("localhost", 5432, "bobayka", "12345", "TinkoffProj")
	if err != nil {
		log.Fatalf("pginit: %s", err)
	}
	var auth AuthHandler
	r := auth.Routes()
	r.Mount("/api/v1", r)
	log.Fatal(http.ListenAndServe(":5000", r))
}
