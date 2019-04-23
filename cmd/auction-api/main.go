package main

import (
	_ "github.com/lib/pq"
	"gitlab.com/bobayka/courseproject/cmd/Protobuf"
	"gitlab.com/bobayka/courseproject/internal/postgres"
	"gitlab.com/bobayka/courseproject/internal/postgres/storage"
	"gitlab.com/bobayka/courseproject/internal/services"
	"google.golang.org/grpc"
	"log"
	"net"
)

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

	listen, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("can't listen on port: %v", err)
	}
	serv := services.LotService{LotStmtsStorage: storage.Lots}
	s := grpc.NewServer()
	lotspb.RegisterLotsServiceServer(s, &serv)
	if err := s.Serve(listen); err != nil {
		log.Fatalf("can't register service server: %v", err)
	}
}
