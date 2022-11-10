package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/alexyslozada/shorturl/model"
)

const (
	shortURLTable = "short_urls"
)

const (
	sqlShortURLInsert = `INSERT INTO ` + shortURLTable + ` (id, short, redirect_to, description, created_at) 
				VALUES ($1, $2, $3, $4, $5)`
	sqlShortURLUpdate = `UPDATE ` + shortURLTable + `
				SET short = $1, redirect_to = $2, description = $3, updated_at = $4
				WHERE id = $5`
	sqlShortURLIncrement    = `UPDATE ` + shortURLTable + ` SET times = times + 1 WHERE id = $1`
	sqlShortURLDelete       = `DELETE FROM ` + shortURLTable + ` WHERE id = $1`
	sqlShortURLQuery        = `SELECT id, short, redirect_to, description, times, created_at, updated_at FROM ` + shortURLTable
	sqlGetAll               = sqlShortURLQuery + ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	sqlShortURLQueryByShort = sqlShortURLQuery + ` WHERE short = $1`
)

type ShortURL struct {
	db *pgxpool.Pool
}

func NewShortURL(db *pgxpool.Pool) ShortURL {
	return ShortURL{db: db}
}

func (s ShortURL) Create(m *model.ShortURL) error {
	_, err := s.db.Exec(
		context.TODO(),
		sqlShortURLInsert,
		m.ID,
		m.Short,
		m.RedirectTo,
		m.Description,
		m.CreatedAt,
	)

	return err
}

func (s ShortURL) Update(m *model.ShortURL) error {
	_, err := s.db.Exec(
		context.TODO(),
		sqlShortURLUpdate,
		m.Short,
		m.RedirectTo,
		m.Description,
		m.UpdatedAt,
		m.ID,
	)

	return err
}

// IncrementTimes never will have a race condition b/c postgres has a `READ COMMITTED` isolation level by default.
// @see https://www.postgresql.org/docs/current/transaction-iso.html#XACT-READ-COMMITTED
func (s ShortURL) IncrementTimes(ID uuid.UUID) error {
	_, err := s.db.Exec(
		context.TODO(),
		sqlShortURLIncrement,
		ID,
	)

	return err
}

func (s ShortURL) Delete(ID uuid.UUID) error {
	_, err := s.db.Exec(
		context.TODO(),
		sqlShortURLDelete,
		ID,
	)

	return err
}

func (s ShortURL) ByShort(short string) (model.ShortURL, error) {
	row := s.db.QueryRow(
		context.TODO(),
		sqlShortURLQueryByShort,
		short,
	)

	return s.scan(row)
}

func (s ShortURL) All(limit, offset int) (model.ShortURLs, error) {
	rows, err := s.db.Query(
		context.TODO(),
		sqlGetAll,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms model.ShortURLs
	for rows.Next() {
		m, err := s.scan(rows)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	return ms, nil
}

func (s ShortURL) scan(row pgx.Row) (model.ShortURL, error) {
	m := model.ShortURL{}

	updatedAtNull := sql.NullInt64{}

	err := row.Scan(
		&m.ID,
		&m.Short,
		&m.RedirectTo,
		&m.Description,
		&m.Times,
		&m.CreatedAt,
		&updatedAtNull,
	)

	m.UpdatedAt = updatedAtNull.Int64

	return m, err
}
