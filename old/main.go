package main

import (
	"context"
	"fmt"
	"github.com/ashikhman/scraper/old/pkg/db"
	"github.com/ashikhman/scraper/old/pkg/importer"
	"github.com/ashikhman/scraper/old/pkg/scraper"
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
	for worker := 0; worker < 100; worker++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := range tasks {
				loader.Load(i)
			}
		}()
	}

	// https://individualka-online.com/wp-content/uploads_n1/bigp36310_1.jpg
	for i := 32000; i < 36310; i++ {
		for j := 1; j <= 20; j++ {
			tasks <- fmt.Sprintf("https://individualka-online.com/wp-content/uploads_n1/bigp%d_%d.jpg", i, j)
		}
	}

	// https://intimka.nl/Persons_deleted/260/Big/260000/01.jpg
	//for i := 306100; i < 370000; i++ {
	//	for f := 1; f < 21; f++ {
	//		firstThree := i / 1000
	//		tasks <- fmt.Sprintf("https://intimka.nl/Persons_deleted/%d/Big/%d/%02d.jpg", firstThree, i, f)
	//	}
	//}

	//for i := 184140; i < 186961; i++ {
	//	tasks <- fmt.Sprintf("http://b.intimdialog.net/img/upload/0%d.jpg", i)
	//}
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
