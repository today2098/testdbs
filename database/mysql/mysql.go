package mysql

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mmy "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/today2098/testdbs/database"
)

func init() {
	database.Register("mysql", &Mysql{})
}

type Mysql struct{}

var _ database.Driver = (*Mysql)(nil)

func (d *Mysql) Open(dsn, dbName string) (*sql.DB, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	// set database name
	cfg.DBName = dbName

	cfg.MultiStatements = true // !

	// create a new *sql.DB to connect a test database
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *Mysql) NewMigrate(db *sql.DB, sourceUrl string) (*migrate.Migrate, error) {
	driver, err := mmy.WithInstance(db, &mmy.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance(sourceUrl, "mysql", driver)
	if err != nil {
		return nil, err
	}
	return m, nil
}
