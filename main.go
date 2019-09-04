package main

import (
	"context"
	"fmt"
	"github.com/ashikhman/scraper/pkg/db"
	"github.com/ashikhman/scraper/pkg/importer"
	"github.com/ashikhman/scraper/pkg/scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/rand"
	"os"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	db.Connect()

	loader := scraper.NewLoader()

	loader.Load("http://en.tool.ws.pho.to/img/howto_images/face_creating/1/source.jpg")

	tasks := make(chan string)
	var wg sync.WaitGroup
	for worker := 0; worker < 1; worker++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := range tasks {
				loader.Load(i)
			}
		}()
	}

	for i := 155900; i < 165900; i++ {
		tasks <- fmt.Sprintf("%d.jpg", i)
	}
	close(tasks)
	wg.Wait()

	//fillQueue()
	//
	//Runner := scraper.NewRunner(context.Background(), db.Pool())
	//defer Runner.Release()
	//
	//err := Runner.Run()
	//if err != nil {
	//	log.Fatal().Err(err).Msg("Runner has failed")
	//}
	//
	//log.Info().Msg("Done")
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
