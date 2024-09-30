package main

import (
	"database/sql"
	"github.com/WooDMaNbtw/BankApp/utils"

	"github.com/WooDMaNbtw/BankApp/api"
	db "github.com/WooDMaNbtw/BankApp/db/sqlc"
	_ "github.com/lib/pq"
	"log"
)

// Entrypoint for HTTP starting server
func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
