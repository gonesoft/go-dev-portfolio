package testutils

import (
	"database/sql"
	"gonesoft/go-dev-portfolio/internal/db"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var TestDB *sql.DB

func TestMain(m *testing.M) {
	log.Println("🔧 Setting up test database...")

	TestDB = db.Connect() // Connect using test env vars

	if err := TestDB.Ping(); err != nil {
		log.Fatalf("❌ Cannot connect to test DB: %v", err)
	}

	log.Println("✅ Test database connected")

	code := m.Run()

	log.Println("🧹 Cleaning up test database...")
	TestDB.Close()

	os.Exit(code)
}
