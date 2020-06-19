package main

import (
	"github.com/ashikhman/scraper/pkg/entity"
	"github.com/dgraph-io/badger/v2"
	"log"
)

func main() {
	var path = "/home/vadim/projects/github.com/ashikhman/scraper/db"
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var aga = entity.NewSource()
	aga.Domain = "DSADS"
}
