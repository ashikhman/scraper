package main

import (
	"context"
	"fmt"
	"github.com/ashikhman/scraper/pkg/db"
	"github.com/ashikhman/scraper/pkg/importer"
	"github.com/ashikhman/scraper/pkg/scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	db.Connect()

	fillQueue()

	Runner := scraper.NewRunner(context.Background(), db.Pool())
	defer Runner.Release()

	err := Runner.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Runner has failed")
	}

	log.Info().Msg("Done")
}

func fillQueue() {
	_, err := db.Pool().Exec(context.Background(), "TRUNCATE TABLE queue_record")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to truncate queue")
	}

	template := "http://localhost:55002/images/%d.jpg"
	Importer := importer.New("localhost:55002")

	for i := 1; i <= 4; i++ {
		Importer.AddUrl(fmt.Sprintf(template, i))
	}

	if err := Importer.Flush(); err != nil {
		log.Fatal().Err(err).Msg("Failed to flush Importer")
	}
}
