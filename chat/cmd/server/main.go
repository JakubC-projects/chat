package main

import (
	"context"
	"log"
	"os"

	"github.com/JakubC-projects/chat/chat/postgres"
	"github.com/JakubC-projects/chat/chat/pubsub"
	"github.com/JakubC-projects/chat/chat/server"
)

var (
	port         = getEnv("PORT", "3001")
	publicDir    = getEnv("PUBLIC_DIR", "./public")
	dbConnString = getEnv("DB_CONN_STRING", "dbname=chat")
)

func main() {
	db := postgres.NewDb(dbConnString)
	dbEvents, err := db.Subscribe(context.Background())
	if err != nil {
		panic(err)
	}
	pubsub := pubsub.New(dbEvents)
	go pubsub.Run()
	srv := server.New(port, db, pubsub, publicDir)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}

func getEnv(name, defaultValue string) string {
	val, ok := os.LookupEnv(name)
	if ok {
		return val
	}
	return defaultValue
}
