package scraper

import (
	"context"
	"errors"
	"github.com/ashikhman/scraper/pkg/db"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/proxy"
	"net"
	"time"
)

type Runner struct {
	ctx  context.Context
	pool *pgxpool.Pool
	tx   *pgxpool.Tx
}

func NewRunner(ctx context.Context, pool *pgxpool.Pool) *Runner {
	return &Runner{
		ctx:  ctx,
		pool: pool,
	}
}

func (r *Runner) Run() error {
	var err error

	r.tx, err = r.pool.Begin(r.ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		txErr := r.tx.Err()
		if txErr != nil {
			err := r.tx.Rollback(r.ctx)
			if err != nil {
				log.Error().Err(txErr).Err(err).Msg("Failed to rollback runner's transaction")
			}

			log.Info().Msg("Runner's transaction committed")
		} else {
			err := r.tx.Commit(r.ctx)
			if err != nil {
				log.Error().Err(err).Msg("Failed to commit runner's transaction")
			}

			log.Info().Msg("Runner's transaction committed")
		}
	}()

	rows, err := r.tx.Query(r.ctx, "SELECT domain_id, url FROM queue_record LIMIT 100 FOR UPDATE SKIP LOCKED")
	if err != nil {
		return err
	}
	defer rows.Close()

	proxyServer, err := db.FindRandomProxyServer()
	if err != nil {
		return err
	}
	if !proxyServer.Auth.Defined() {
		return errors.New("no auth defined for proxy server")
	}

	log.Info().Str("address", proxyServer.Socks5Address()).Str("username", proxyServer.Auth.Username.String).
		Msg("Using proxy server")
	proxyDial := func(addr string) (net.Conn, error) {
		dialer, err := proxy.SOCKS5("tcp", proxyServer.Socks5Address(), &proxy.Auth{
			User:     proxyServer.Auth.Username.String,
			Password: proxyServer.Auth.Password.String,
		}, proxy.Direct)
		if err != nil {
			return nil, err
		}
		return dialer.Dial("tcp", addr)
	}
	_ = proxyDial

	http := &fasthttp.Client{
		//Dial:         proxyDial,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	var (
		domainId int
		url      string
	)
	for rows.Next() {
		err = rows.Scan(&domainId, &url)
		if err != nil {
			return err
		}

		log.Debug().Str("url", url).Msg("Processing url")

		request.SetRequestURI(url)
		err = http.Do(request, response)
		if err != nil {
			return err
		}

	}

	return nil
}

func (r *Runner) Release() {

}
