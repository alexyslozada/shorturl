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
	historyTable = "histories"
)

const (
	sqlHistoryInsert                    = `INSERT INTO ` + historyTable + ` (id, short_url_id, created_at) VALUES ($1, $2, $3)`
	sqlHistoryQuery                     = `SELECT id, short_url_id, created_at, updated_at FROM ` + historyTable
	sqlHistoryQueryByShortURLID         = sqlHistoryQuery + ` WHERE short_url_id = $1`
	sqlHistoryQueryByShortURLIDAndDates = sqlHistoryQuery + ` WHERE short_url_id = $1 AND created_at BETWEEN $2 AND $3`
)

type History struct {
	db *pgxpool.Pool
}

func NewHistory(db *pgxpool.Pool) History {
	return History{db: db}
}

func (h History) Create(m *model.History) error {
	_, err := h.db.Exec(
		context.TODO(),
		sqlHistoryInsert,
		m.ID,
		m.ShortURLID,
		m.CreatedAt,
	)

	return err
}

func (h History) ByShortURLID(ID uuid.UUID) (model.Histories, error) {
	rows, err := h.db.Query(
		context.TODO(),
		sqlHistoryQueryByShortURLID,
		ID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return h.Histories(rows)
}

func (h History) ByShortURLIDAndDates(ID uuid.UUID, from, to int64) (model.Histories, error) {
	rows, err := h.db.Query(
		context.TODO(),
		sqlHistoryQueryByShortURLIDAndDates,
		ID,
		from,
		to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return h.Histories(rows)
}

func (h History) All() (model.Histories, error) {
	rows, err := h.db.Query(
		context.TODO(),
		sqlHistoryQuery,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return h.Histories(rows)
}

func (h History) scan(row pgx.Row) (model.History, error) {
	m := model.History{}

	updatedAtNull := sql.NullInt64{}

	err := row.Scan(
		&m.ID,
		&m.ShortURLID,
		&m.CreatedAt,
		&updatedAtNull,
	)

	m.UpdatedAt = updatedAtNull.Int64

	return m, err
}

func (h History) Histories(rows pgx.Rows) (model.Histories, error) {
	var ms model.Histories
	for rows.Next() {
		m, err := h.scan(rows)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	return ms, nil
}
