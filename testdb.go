package testdbs

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type TestDatabase struct {
	dbName string
	db     *sql.DB
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

// Drop closes and drops the test database.
func (td *TestDatabase) Drop() error {
	return td.par.Drop(td)
}
