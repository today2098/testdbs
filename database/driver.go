package database

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/golang-migrate/migrate/v4"
)

var drivers sync.Map

type Driver interface {
	Open(dsn, dbName string) (*sql.DB, error)
	NewMigrate(db *sql.DB, sourceUrl string) (*migrate.Migrate, error)
}

// Open returns a new *sql.DB.
func Open(driverName, dsn, dbName string) (*sql.DB, error) {
	d, exists := drivers.Load(driverName)
	if !exists {
		return nil, fmt.Errorf("source driver: unknown driver '%s' (forgotten import?)", driverName)
	}

	return d.(Driver).Open(dsn, dbName)
}

// NewMigrate returns a new *migrate.Migrate.
func NewMigrate(driverName string, db *sql.DB, sourceUrl string) (*migrate.Migrate, error) {
	d, exists := drivers.Load(driverName)
	if !exists {
		return nil, fmt.Errorf("source driver: unknown driver '%s' (forgotten import?)", driverName)
	}

	return d.(Driver).NewMigrate(db, sourceUrl)
}

// Register globally registers a driver.
func Register(name string, driver Driver) {
	if driver == nil {
		panic("Register driver is nil")
	}
	if _, exists := drivers.Load(name); exists {
		panic("Register called twice for driver " + name)
	}
	drivers.Store(name, driver)
}
