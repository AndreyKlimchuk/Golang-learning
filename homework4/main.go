package main

import (
	"log"

	"github.com/AndreyKlimchuk/golang-learning/homework4/db"

	"github.com/AndreyKlimchuk/golang-learning/homework4/api"
	"github.com/AndreyKlimchuk/golang-learning/homework4/logger"
)

func main() {
	if err := logger.InitZap(); err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	conn := "user=gorello password=12345 host=localhost port=5432 dbname=gorello pool_max_conns=10"
	if err := db.Init(conn); err != nil {
		log.Fatalf("can't initialize db: %v", err)
	}
	defer func() {
		if err := logger.Zap.Sync(); err != nil {
			log.Print("can't sync zap logger")
		}
	}()
	api.StartHttpServer()
}
