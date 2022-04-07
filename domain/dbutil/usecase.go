package dbutil

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBUtil struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) DBUtil {
	return DBUtil{db: db}
}

func (dbu DBUtil) Tx() (pgx.Tx, error) {
	return dbu.db.Begin(context.TODO())
}
