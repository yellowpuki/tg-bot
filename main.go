package main

import (
	"context"
	"flag"
	"log"

	tgClient "github.com/yellowpuki/tg-bot/clients/telegram"
	event_consumer "github.com/yellowpuki/tg-bot/consumer/event-consumer"
	"github.com/yellowpuki/tg-bot/events/telegram"
	"github.com/yellowpuki/tg-bot/storage/sqlite"
)

const (
	tgBotHost          = "api.telegram.org"
	filesStoragePath   = "data/files"
	sqlliteStoragePath = "data/sqlite/storage.db"
	batchSize          = 100
)

func main() {
	//s := files.New(storagePath)

	s, err := sqlite.New(sqlliteStoragePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage", err)
	}

	eventProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}

// mustToken ...
func mustToken() string {
	token := flag.String(
		"t",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
