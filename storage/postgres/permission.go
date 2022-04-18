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
	permissionTable = "permissions"
)

const (
	sqlPermissionInsert = `INSERT INTO ` + permissionTable + ` 
				(id, user_id, can_create, can_update, can_delete, can_select, is_admin, created_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	sqlPermissionUpdate = `UPDATE ` + permissionTable + `
				SET can_create = $1, can_update= $2, can_delete = $3, 
					can_select = $4, is_admin = $5, updated_at = $6
				WHERE id = $7`
	sqlPermissionDelete = `DELETE FROM ` + permissionTable + ` WHERE id = $1`
	sqlPermissionQuery  = `SELECT id, user_id, can_create, can_update, can_delete, 
					can_select, is_admin, created_at, updated_at
				FROM ` + permissionTable
	sqlPermissionQueryByUserID = sqlPermissionQuery + ` WHERE user_id = $1`
)

type Permission struct {
	db *pgxpool.Pool
}

func NewPermission(db *pgxpool.Pool) Permission {
	return Permission{db: db}
}

func (p Permission) Create(m *model.Permission) error {
	_, err := p.db.Exec(
		context.TODO(),
		sqlPermissionInsert,
		m.ID,
		m.UserID,
		m.CanCreate,
		m.CanUpdate,
		m.CanDelete,
		m.CanSelect,
		m.IsAdmin,
		m.CreatedAt,
	)

	return err
}

func (p Permission) Update(m *model.Permission) error {
	_, err := p.db.Exec(
		context.TODO(),
		sqlPermissionUpdate,
		m.CanCreate,
		m.CanUpdate,
		m.CanDelete,
		m.CanSelect,
		m.IsAdmin,
		m.UpdatedAt,
		m.ID,
	)

	return err
}

func (p Permission) Delete(ID uuid.UUID) error {
	_, err := p.db.Exec(
		context.TODO(),
		sqlPermissionDelete,
		ID,
	)

	return err
}

func (p Permission) ByUserID(ID uuid.UUID) (model.Permission, error) {
	row := p.db.QueryRow(
		context.TODO(),
		sqlPermissionQueryByUserID,
		ID,
	)

	return p.scan(row)
}

func (p Permission) All() (model.Permissions, error) {
	rows, err := p.db.Query(
		context.TODO(),
		sqlPermissionQuery,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms model.Permissions
	for rows.Next() {
		m, err := p.scan(rows)
		if err != nil {
			return nil, err
		}

		ms = append(ms, m)
	}

	return ms, nil
}

func (p Permission) scan(row pgx.Row) (model.Permission, error) {
	m := model.Permission{}
	updateNull := sql.NullInt64{}

	err := row.Scan(
		&m.ID,
		&m.UserID,
		&m.CanCreate,
		&m.CanUpdate,
		&m.CanDelete,
		&m.CanSelect,
		&m.IsAdmin,
		&m.CreatedAt,
		&updateNull,
	)
	m.UpdatedAt = updateNull.Int64

	return m, err
}
