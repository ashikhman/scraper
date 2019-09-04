package scraper

import (
	"bytes"
	"github.com/ashikhman/scraper/pkg/db"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type Loader struct {
	client      *fasthttp.Client
	storagePath string
	storage     *Storage
}

func NewLoader() *Loader {
	client := &fasthttp.Client{
		Dial:         proxyDialFunc(),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	//for i := 0; i < 10000; i++ {
	//	err := os.Mkdir(fmt.Sprintf("/storage/%d", i), 0644)
	//	if err != nil {
	//		log.Fatal().Err(err).Msg("Failed to create a directory")
	//	}
	//}

	return &Loader{
		client:      client,
		storagePath: strings.TrimSuffix("/storage", "/"),
		storage:     NewStorage(),
	}
}

func (l *Loader) Load(address string) (item *Item) {
	log.Info().Str("address", address).Msg("Loading address")

	item = l.storage.Get(address)
	if item != nil {
		return item
	}

	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	request.SetRequestURI(address)
	err := l.client.Do(request, response)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to execute HTTP request")
	}

	if response.StatusCode() < 200 || response.StatusCode() >= 500 {
		log.Fatal().Str("url", address).Int("status_code", response.StatusCode()).Msg("Invalid response code")
	}

	fileName, err := uuid.NewRandom()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create random UUID")
	}

	path := bytes.Buffer{}
	path.WriteString(l.storagePath)
	path.WriteRune('/')
	path.WriteString(strconv.Itoa(rand.Intn(10000)))
	path.WriteRune('/')
	path.WriteString(fileName.String())

	err = ioutil.WriteFile(path.String(), response.Body(), 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to write to file")
	}

	item = &Item{
		Url:        address,
		ContentUri: path.String(),
		StatusCode: response.StatusCode(),
	}
	l.storage.Save(item)

	return item
}

func proxyDialFunc() fasthttp.DialFunc {
	proxyServer, err := db.FindRandomProxyServer()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to acquire proxy server")
	}
	if !proxyServer.Auth.Defined() {
		log.Fatal().Int("id", proxyServer.ID).Msg("No auth defined for proxy server")
	}

	log.Info().Str("address", proxyServer.Socks5Address()).Msg("Using proxy server")

	return func(addr string) (net.Conn, error) {
		dialer, err := proxy.SOCKS5("tcp", proxyServer.Socks5Address(), &proxy.Auth{
			User:     proxyServer.Auth.Username.String,
			Password: proxyServer.Auth.Password.String,
		}, proxy.Direct)
		if err != nil {
			return nil, err
		}
		return dialer.Dial("tcp", addr)
	}
}
