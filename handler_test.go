package testdbs_test

import "github.com/today2098/testdbs"

func ExampleNewHandler() {
	dsn := "user:password@tcp(localhost:3306)/?multiStatements=true"
	sourceUrl := "file:///home/path/to/your/migrations"
	h := testdbs.NewHandler("mysql", dsn, sourceUrl)
	h.Connect()
	defer h.Close()

	td, _ := h.CreateAndMigrate()
	defer td.Drop()

	dbx := td.DBx()
	dbx.Exec("INSERT INTO tables (id) VALUES (1)")
}
