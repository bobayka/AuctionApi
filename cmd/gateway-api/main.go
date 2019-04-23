package main

import (
	_ "github.com/lib/pq"
	"gitlab.com/bobayka/courseproject/cmd/Protobuf"
	"gitlab.com/bobayka/courseproject/cmd/gateway-api/general"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/postgres/storage"
	"google.golang.org/grpc"
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

	storage, err := storage.NewStorage(db)
	if err != nil {
		log.Fatalf("error in creation usersstorage: %s", err)
	}
	defer storage.Close()

	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can't connect to server: %v", err)
	}
	defer conn.Close()
	client := lotspb.NewLotsServiceClient(conn)

	auth := gatewayApi.NewGatewayApi(storage, client)
	r := auth.Routes()
	go func() {
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()
	log.Fatal(http.ListenAndServe(":5000", r))
}
