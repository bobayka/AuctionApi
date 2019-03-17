package main

import (
	"log"
	"net/http"
)

func main() {
	var auth AuthHandler
	r := auth.Routes()
	r.Mount("/api/v1", r)
	log.Fatal(http.ListenAndServe(":5000", r))
}
