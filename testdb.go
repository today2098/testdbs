package testdbs

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
)

// TestDatabase contains some helper methods to handle a test database.
type TestDatabase struct {
	dbName string
	db     *sql.DB
	m      *migrate.Migrate
	par    *Handler
}

// DB returns *sql.DB.
func (td *TestDatabase) DB() *sql.DB {
	return td.db
}

// DBx returns *sqlx.DB instead for *sql.DB.
func (td *TestDatabase) DBx() *sqlx.DB {
	return sqlx.NewDb(td.DB(), "mysql")
}

// Migrate returns *migrate.Migrate.
func (td *TestDatabase) Migrate() *migrate.Migrate {
	return td.m
}

// Drop closes and drops the test database.
func (td *TestDatabase) Drop() error {
	return td.par.Drop(td)
}
