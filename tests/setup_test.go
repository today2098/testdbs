package tests_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/today2098/testdbs"
	_ "github.com/today2098/testdbs/database/mysql"
)

var h *testdbs.Handler // NOTE: h will be overwrited by TestMain().

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("fatal: %v", err)
	}

	cfg, err := mysql.ParseDSN(os.Getenv("DSN_TEST"))
	if err != nil {
		log.Fatalf("fatal: %v", err)
	}
	cfg.Loc = time.UTC
	cfg.ParseTime = true
	cfg.MultiStatements = true

	h = testdbs.NewHandler("mysql", os.Getenv("DSN_TEST"), "file://./migrations")
	if err := h.Connect(); err != nil {
		log.Fatalf("fatal: %v", err)
	}
	defer h.Close() // NOTE: Do not forget to call Close().

	m.Run()
}
