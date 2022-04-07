package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/alexyslozada/shorturl/model"
)

const (
	userTable = "users"
)

const (
	sqlUserInsert       = `INSERT INTO ` + userTable + ` (id, email, password, full_name, active, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	sqlUserDelete       = `DELETE FROM ` + userTable + ` WHERE id = $1`
	sqlUserQuery        = `SELECT id, email, password, full_name, active, created_at, updated_at FROM ` + userTable
	sqlUserQueryByEmail = sqlUserQuery + ` WHERE email = $1`
)

type User struct {
	db *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) User {
	return User{db: db}
}

func (u User) Create(m *model.User) error {
	_, err := u.db.Exec(
		context.TODO(),
		sqlUserInsert,
		m.ID,
		m.Email,
		m.Password,
		m.FullName,
		m.Active,
		m.CreatedAt,
	)

	return err
}

func (u User) Delete(ID uuid.UUID) error {
	_, err := u.db.Exec(
		context.TODO(),
		sqlUserDelete,
		ID,
	)

	return err
}

func (u User) ByEmail(email string) (model.User, error) {
	row := u.db.QueryRow(
		context.TODO(),
		sqlUserQueryByEmail,
		email,
	)

	return u.scan(row)
}

func (u User) All() (model.Users, error) {
	rows, err := u.db.Query(
		context.TODO(),
		sqlUserQuery,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms model.Users
	for rows.Next() {
		m, err := u.scan(rows)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	return ms, nil
}

func (u User) scan(row pgx.Row) (model.User, error) {
	m := model.User{}
	err := row.Scan(
		&m.ID,
		&m.Email,
		&m.Password,
		&m.FullName,
		&m.Active,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	return m, err
}
