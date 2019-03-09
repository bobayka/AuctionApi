package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

func main() {
	var auth AuthHandler
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/signup", auth.RegistrationHandler)
		r.Post("/signin", auth.AuthorizationHandler)
		r.Put("/users/0", auth.UpdateHandler)
	})
	log.Fatal(http.ListenAndServe(":5000", r))

}
