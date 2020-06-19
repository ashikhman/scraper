package main

import (
	"github.com/dgraph-io/badger"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to open Badger db")
		return
	}
	defer (func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Badger db")
		}
	})()

	time.Sleep(20 * time.Second)

	log.Info().Msg("DONE")

	//wb := db.NewWriteBatch()
	//defer wb.Cancel()
	//
	//for i := 0; i < 10000000; i++ {
	//	err = wb.Set([]byte(RandStringRunes(500)), []byte(RandStringRunes(500)))
	//}
	//
	//err = wb.Flush()
	//if err != nil {
	//	log.Error().Err(err).Msg("Failed to flush batch job")
	//}
}
