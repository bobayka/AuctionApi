package main

import (
	"log"
	"net/http"
)

func main() {
	var auth AuthHandler
	r := auth.Routes()
	log.Fatal(http.ListenAndServe(":5000", r))

}
