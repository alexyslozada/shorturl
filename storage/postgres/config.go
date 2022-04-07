package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	AppName = "go-short-url"
)

const (
	MinConns = 3
	MaxConns = 100
)

func New(user, pass, host, port, dbName, sslMode string, minConns, maxConns int32) (*pgxpool.Pool, error) {
	connString := makeURL(user, pass, host, port, dbName, sslMode)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	config.MinConns = MinConns
	if minConns > 0 && minConns < MinConns {
		config.MinConns = minConns
	}

	config.MaxConns = MaxConns
	if maxConns > 0 && maxConns < MaxConns {
		config.MaxConns = maxConns
	}

	config.ConnConfig.RuntimeParams["application_name"] = AppName

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, err
}

func makeURL(user, pass, host, port, dbName, sslMode string) string {
	return fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=%s", user, pass, host, port, dbName, sslMode)
}
