package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"strings"
)

var pool *pgxpool.Pool

func Connect() {
	databaseUrl := "postgres://scraper:scraper@localhost:55001/scraper"
	Pool, err := pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database.")
	}
	pool = Pool

	log.Info().Msg("Database connection has established.")
}

func Pool() *pgxpool.Pool {
	return pool
}

func Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	return pool.Acquire(ctx)
}

func Close() {
	pool.Close()
}

func quoteString(str string) string {
	return "'" + strings.Replace(str, "'", "''", -1) + "'"
}
