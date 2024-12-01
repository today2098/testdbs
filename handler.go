package testdbs

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mmy "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/oklog/ulid/v2"
)

type Handler struct {
	cfg       *mysql.Config
	db        *sql.DB
	sourceUrl string
	prefix    string
	children  sync.Map
}

// NewHandler returns an object of Handler.
// NOTE: `sourceUrl` have to be represented by "file://[path to file]"
func NewHandler(cfg *mysql.Config, sourceUrl string, prefix string) *Handler {
	return &Handler{
		cfg:       cfg,
		sourceUrl: sourceUrl,
		prefix:    prefix,
		children:  sync.Map{},
	}
}

// Connect connects to a database and verify with a ping.
func (tdh *Handler) Connect() error {
	var err error
	if tdh.db, err = sql.Open("mysql", tdh.cfg.FormatDSN()); err != nil {
		return err
	}
	return tdh.db.Ping()
}

// Create creates a new test DB and returns a *TestDatabase.
func (tdh *Handler) Create() (*TestDatabase, error) {
	var err error

	// create a new test DB
	dbName := tdh.prefix + "_" + ulid.Make().String()
	if _, err := tdh.db.Exec("CREATE DATABASE " + dbName); err != nil {
		return nil, err
	}

	// create a new *TestDatabase
	child := &TestDatabase{
		dbName: dbName,
		par:    tdh,
	}
	tdh.children.Store(child, struct{}{})

	// connect to a new test DB
	cfg := *tdh.cfg
	cfg.DBName = dbName
	cfg.MultiStatements = true // !
	if child.db, err = sql.Open("mysql", cfg.FormatDSN()); err != nil {
		tdh.Drop(child)
		return nil, err
	}

	// migration
	driver, err := mmy.WithInstance(child.db, &mmy.Config{})
	if err != nil {
		tdh.Drop(child)
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(tdh.sourceUrl, "mysql", driver)
	if err != nil {
		tdh.Drop(child)
		return nil, err
	}
	if err := m.Up(); err != nil {
		tdh.Drop(child)
		return nil, err
	}

	return child, nil
}

// Drop closes and drops a test database.
func (tdh *Handler) Drop(child *TestDatabase) error {
	if err := child.db.Close(); err != nil {
		return err
	}
	if _, err := tdh.db.Exec("DROP DATABASE " + child.dbName); err != nil {
		return err
	}
	tdh.children.Delete(child)
	return nil
}

// AllDrop drops all test databases.
func (tdh *Handler) AllDrop() error {
	var errs error
	tdh.children.Range(func(child, _ any) bool {
		if err := child.(*TestDatabase).Drop(); err != nil {
			errs = errors.Join(errs, err)
		}
		return true
	})
	return errs
}

// Close drops all test databases and close the main database.
func (tdh *Handler) Close() error {
	if err := tdh.AllDrop(); err != nil {
		return err
	}
	return tdh.db.Close()
}
