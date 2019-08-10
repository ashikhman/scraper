package db

import (
	"bytes"
	"context"
	"database/sql"
	"strconv"
)

type ProxyServer struct {
	ID         int
	Host       string
	Socks5Port int
	Auth       ProxyServerAuth
}

func (ps *ProxyServer) Socks5Address() string {
	var buffer bytes.Buffer

	buffer.WriteString(ps.Host)
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(ps.Socks5Port))

	return buffer.String()
}

type ProxyServerAuth struct {
	Username sql.NullString
	Password sql.NullString
}

func (psa ProxyServerAuth) Defined() bool {
	return psa.Username.Valid
}

func FindRandomProxyServer() (*ProxyServer, error) {
	query := `
SELECT    ps.id, ps.host, ps.socks5_port, psa.username, psa.password
FROM      proxy_server ps
LEFT JOIN proxy_server_auth psa ON psa.id = ps.auth_id
LIMIT     1
`
	proxyServer := &ProxyServer{}

	err := pool.QueryRow(context.Background(), query).Scan(
		&proxyServer.ID,
		&proxyServer.Host,
		&proxyServer.Socks5Port,
		&proxyServer.Auth.Username,
		&proxyServer.Auth.Password,
	)
	if err != nil {
		return nil, err
	}

	return proxyServer, nil
}
