package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/proyuen/simple-bank/api"
	db "github.com/proyuen/simple-bank/db/sqlc"
	"github.com/proyuen/simple-bank/gapi"
	"github.com/proyuen/simple-bank/pb"
	"github.com/proyuen/simple-bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.Newstore(conn)
	runGrpcServer(config, store)

}
func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimplebankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServer_Address)
	if err != nil {
		log.Fatal("cannot create listener")
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gPRC server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServer_Address)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

}
