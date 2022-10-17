package main

import (
	"database/sql"
	"log"
	"simpleauth/api"
	db "simpleauth/db/sqlc"

	"simpleauth/util"

	_ "github.com/lib/pq"
	"github.com/wneessen/go-mail"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	mailer := mail.NewMsg()
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store, mailer)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
