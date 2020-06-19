package scraper

import (
	"context"
	"github.com/ashikhman/scraper/old/pkg/db"
	"github.com/jackc/pgtype"
	"github.com/rs/zerolog/log"
	"sync"
)

type Storage struct {
	lock sync.RWMutex
	ctx  context.Context
}

type Item struct {
	Url        string
	ContentUri pgtype.Text
	StatusCode int
}

func NewStorage() *Storage {
	return &Storage{
		ctx: context.Background(),
	}
}

func (s *Storage) Get(url string) (item *Item) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	query := `
SELECT url, content_uri, status_code 
FROM   storage 
WHERE  url = $1
`

	rows, err := db.Pool().Query(s.ctx, query, url)
	if err != nil {
		log.Fatal().Err(err).Str("url", url).Msg("Failed to fetch storage record")
	}
	defer rows.Close()

	for rows.Next() {
		item = &Item{}
		err := rows.Scan(&item.Url, &item.ContentUri, &item.StatusCode)
		if err != nil {
			log.Fatal().Err(err).Msg("Scan() to Item has failed")
		}

		return item
	}

	return nil
}

func (s *Storage) Save(item *Item) {
	s.lock.Lock()
	defer s.lock.Unlock()

	query := `
INSERT INTO storage (url, content_uri, status_code)
VALUES ($1, $2, $3)
`
	_, err := db.Pool().Exec(s.ctx, query, item.Url, &item.ContentUri, item.StatusCode)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to insert storage record")
	}
}
