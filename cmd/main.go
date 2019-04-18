package main

import (
	_ "github.com/lib/pq"
	"gitlab.com/bobayka/courseproject/cmd/auth-api"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
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

	postgres.StartDBBackgroundProcesses(StmtsStorage)

	auth := authapi.NewAuthApi(StmtsStorage)
	r := auth.Routes()
	go func() {
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()
	log.Fatal(http.ListenAndServe(":5000", r))
}
