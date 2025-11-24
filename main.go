package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/proyuen/simplebank/api"
	db "github.com/proyuen/simplebank/db/sqlc"
	"github.com/proyuen/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("not connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
