# testdbs

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![GoDoc](https://godoc.org/github.com/Songmu/flextime?status.svg)][godoc]

[license]: https://github.com/today2098/testdbs/blob/main/LICENSE
[godoc]: https://godoc.org/github.com/today2098/testdbs

testdbs enables parallel testing in unit tests that use a real database (MySQL only)."


## Synopsis

```go
import "github.com/today2098/testdbs"

var h *testdbs.Handler // NOTE: h will be overwrited by TestMain()

func TestMain(m *testing.M) {
    // Create a handler to create and drop test databases
    dsn := "user:password@tcp(localhost:3306)/?multiStatements=true"
	sourceUrl := "file:///home/path/to/your/migrations"
    h = testdbs.NewHandler("mysql", dsn, sourceUrl)
    h.Connect()
    defer h.Close() // Drop all test databases after all tests finish

    m.Run()
}

func TestA(t *testing.T) {
    t.Parallel()

    td, _ := h.Create() // Create a test database
    defer td.Drop()     // Drop the test database after TestA
    h.Migration().Up()  // Migrate

    db := td.DB() // Return *sql.DB

    // TODO: implement test-A
}

func TestB(t *testing.T) {
    t.Parallel()

    td, _ := h.CreateAndMigrate() // Create another test database and migrate
    defer td.Drop()               // Drop the test database after TestB

    dbx := td.DBx() // Return *sqlx.DB

    // TODO: implement test-B
}
```


## Description

`testdbs` makes it easy to create separate databases for each test.

Separating databases for each test has the following advantages: records do not conflict between tests, and parallelization of tests becomes easier.

`testdbs` supports the [`sql`](https://pkg.go.dev/database/sql) and [`sqlx`](https://pkg.go.dev/github.com/jmoiron/sqlx) packages as database access interfaces. 
Additionally, it uses the [`golang-migrate`](https://pkg.go.dev/github.com/golang-migrate/migrate/v4) package for database migrations and seeding.


## Installation

```console
$ go get github.com/today2098/testdbs
```
