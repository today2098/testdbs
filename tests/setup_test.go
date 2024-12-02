package tests_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/today2098/testdbs"
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
	cfg.Collation = "utf8mb4_bin"
	cfg.Loc = time.UTC
	cfg.ParseTime = true

	h = testdbs.NewHandler(cfg, "file://./migrations", "testdbs_test")
	if err := h.Connect(); err != nil {
		log.Fatalf("fatal: %v", err)
	}
	defer h.Close() // NOTE: Do not forget to call Close().

	m.Run()
}
