package main

import (
	"Backend/api"
	db "Backend/db/sqlc"
	"Backend/db/utils"
	"database/sql"
	"log"

	
	_ "github.com/lib/pq"
)



func main() {
	config, err := utils.LoadConfig(".")
	if err!=nil{
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Could not connect to db", err)
	}
	store :=db.NewStore(conn)
	server, err :=api.NewServer(config, store)
	if err != nil {
		log.Fatal("Could not create server", err)
	}
	err =server.Start(config.ServerAddress)
	if err!=nil{
		log.Fatal("cannot start server:", err)

	}
}
