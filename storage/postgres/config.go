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
	if minConns > 0 && minConns < MinConns {
		minConns = MinConns
	}
	if maxConns > 0 && maxConns > MaxConns {
		maxConns = MaxConns
	}

	connString := makeURL(user, pass, host, port, dbName, sslMode, minConns, maxConns)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	config.ConnConfig.RuntimeParams["application_name"] = AppName

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return pool, err
}

func makeURL(user, pass, host, port, dbName, sslMode string, minConns, maxConns int32) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_min_conns=%d pool_max_conns=%d",
		user,
		pass,
		host,
		port,
		dbName,
		sslMode,
		minConns,
		maxConns,
	)
}
