package testdbs

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/oklog/ulid/v2"
	"github.com/today2098/testdbs/database"
)

var ErrNilPointer = errors.New("testdbs: nil pointer")

// Handler contains some helper methods to create and drop test databases.
type Handler struct {
	db         *sql.DB
	driverName string
	dsn        string
	sourceUrl  string
	prefix     string
	children   *sync.Map
}

// NewHandler returns an object of Handler.
// NOTE: `sourceUrl` have to be represented by "file://[path to file]"
func NewHandler(driverName, dsn, sourceUrl string) *Handler {
	return &Handler{
		driverName: driverName,
		dsn:        dsn,
		sourceUrl:  sourceUrl,
		prefix:     "testdbs",
		children:   &sync.Map{},
	}
}

// DB returns *sql.DB.
func (h *Handler) DB() *sql.DB {
	return h.db
}

// Connect connects to a database and verify with a ping.
func (h *Handler) Connect() error {
	var err error
	if h.db, err = sql.Open(h.driverName, h.dsn); err != nil {
		return err
	}
	return h.db.Ping()
}

// Create creates a new test DB and returns a *TestDatabase.
func (h *Handler) Create() (*TestDatabase, error) {
	var err error

	// create a new test DB
	dbName := h.prefix + "_" + ulid.Make().String()
	if _, err := h.db.Exec("CREATE DATABASE " + dbName); err != nil {
		return nil, err
	}

	// create a new *TestDatabase
	child := &TestDatabase{
		dbName: dbName,
		par:    h,
	}
	h.children.Store(child, struct{}{})

	// connect to a new test DB
	if child.db, err = database.Open(h.driverName, h.dsn, dbName); err != nil {
		h.Drop(child)
		return nil, err
	}
	if err := child.db.Ping(); err != nil {
		h.Drop(child)
		return nil, err
	}

	// create a new *migrate.Migrate
	if child.m, err = database.NewMigrate(h.driverName, child.db, h.sourceUrl); err != nil {
		h.Drop(child)
		return nil, err
	}

	return child, nil
}

// CreateAndMigrate creates a new test DB, migrates and returns a *TestDatabase.
func (h *Handler) CreateAndMigrate() (*TestDatabase, error) {
	// create a new *TestDatabase
	child, err := h.Create()
	if err != nil {
		return nil, err
	}

	// migrate
	if err := child.Migrate().Up(); err != nil {
		h.Drop(child)
		return nil, err
	}

	return child, nil
}

// Drop closes and drops a test database.
func (h *Handler) Drop(child *TestDatabase) error {
	if child == nil {
		return ErrNilPointer
	}
	if child.db != nil {
		if err := child.db.Close(); err != nil {
			return err
		}
	}
	if _, err := h.db.Exec("DROP DATABASE IF EXISTS " + child.dbName); err != nil {
		return err
	}
	h.children.Delete(child)
	return nil
}

// AllDrop drops all test databases.
func (h *Handler) AllDrop() error {
	var errs error
	h.children.Range(func(child, _ any) bool {
		if err := h.Drop(child.(*TestDatabase)); err != nil {
			errs = errors.Join(errs, err)
		}
		return true
	})
	return errs
}

// Close drops all test databases and close the main database.
func (h *Handler) Close() error {
	if err := h.AllDrop(); err != nil {
		return err
	}
	return h.db.Close()
}
