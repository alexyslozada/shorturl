package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Migrate(db *pgxpool.Pool) error {
	err := createShortURL(db)
	if err != nil {
		return fmt.Errorf("%s %w", "short_url", err)
	}

	err = createHistory(db)
	if err != nil {
		return fmt.Errorf("%s %w", "history", err)
	}

	err = createUser(db)
	if err != nil {
		return fmt.Errorf("%s %w", "user", err)
	}

	err = createPermission(db)
	if err != nil {
		return fmt.Errorf("%s %w", "permission", err)
	}

	return nil
}

func createShortURL(db *pgxpool.Pool) error {
	sql := `CREATE TABLE IF NOT EXISTS ` + shortURLTable + `(
			id UUID NOT NULL,
			short VARCHAR(200) NOT NULL,
			redirect_to VARCHAR(1024) NOT NULL,
			description VARCHAR(1024),
			times INTEGER NOT NULL DEFAULT 0,
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT,
			CONSTRAINT ` + shortURLTable + `_id_pk PRIMARY KEY (id),
			CONSTRAINT ` + shortURLTable + `_short_uk UNIQUE (short)
		)`
	_, err := db.Exec(context.TODO(), sql)

	return err
}

func createHistory(db *pgxpool.Pool) error {
	sql := `
		CREATE TABLE IF NOT EXISTS ` + historyTable + `(
			id UUID NOT NULL,
			short_url_id UUID NOT NULL,
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT,
			CONSTRAINT ` + historyTable + `_id_pk PRIMARY KEY (id),
			CONSTRAINT ` + historyTable + `_short_url_id_fk FOREIGN KEY (short_url_id) REFERENCES ` + shortURLTable + ` (id)
				ON UPDATE RESTRICT ON DELETE RESTRICT
		)`

	_, err := db.Exec(context.TODO(), sql)

	return err
}

func createUser(db *pgxpool.Pool) error {
	sql := `
		CREATE TABLE IF NOT EXISTS ` + userTable + `(
			id UUID NOT NULL,
			email VARCHAR(128) NOT NULL,
			password VARCHAR(256) NOT NULL,
			full_name VARCHAR(256) NOT NULL,
			active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT,
			CONSTRAINT ` + userTable + `_id_pk PRIMARY KEY (id),
			CONSTRAINT ` + userTable + `_email_uk UNIQUE (email)
		)`

	_, err := db.Exec(context.TODO(), sql)

	return err
}

func createPermission(db *pgxpool.Pool) error {
	sql := `
		CREATE TABLE IF NOT EXISTS ` + permissionTable + `(
			id UUID NOT NULL,
			user_id UUID NOT NULL,
			can_create BOOLEAN NOT NULL DEFAULT FALSE,
			can_update BOOLEAN NOT NULL DEFAULT FALSE,
			can_delete BOOLEAN NOT NULL DEFAULT FALSE,
			can_select BOOLEAN NOT NULL DEFAULT FALSE,
			is_admin BOOLEAN NOT NULL DEFAULT FALSE,
			created_at BIGINT NOT NULL DEFAULT 0,
			updated_at BIGINT,
			CONSTRAINT ` + permissionTable + `_id_pk PRIMARY KEY (id),
			CONSTRAINT ` + permissionTable + `_user_id_uk UNIQUE (user_id),
			CONSTRAINT ` + permissionTable + `_user_id_fk FOREIGN KEY (user_id) REFERENCES ` + userTable + ` (id)
				ON UPDATE RESTRICT ON DELETE RESTRICT
		)`

	_, err := db.Exec(context.TODO(), sql)

	return err
}
