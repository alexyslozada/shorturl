package dbutil

import "github.com/jackc/pgx/v4"

type UseCase interface {
	Tx() (pgx.Tx, error)
}
