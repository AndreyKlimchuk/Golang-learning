package main

import (
	"log"

	"github.com/AndreyKlimchuk/golang-learning/homework4/api"
	"github.com/AndreyKlimchuk/golang-learning/homework4/logger"
)

func main() {
	if err := logger.InitZap(); err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err := logger.Zap.Sync(); err != nil {
			log.Print("cannot sync zap logger")
		}
	}()
	api.StartHttpServer()
}
