package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var (
	db   *sql.DB
	once sync.Once
)

// func getEnv() {
// 	cwd, _ := os.Getwd()
// 	root := filepath.Join(cwd, ".env")
// 	if _, err := os.Stat(root); err == nil {
// 		_ = godotenv.Load(root)
// 	} else {
// 		log.Printf("No .env file found in %s", root)
// 	}
// }

// Connect establishes a connection to the database using environment variables.
func Connect() *sql.DB {

	once.Do(func() {

		// Load environment variables look for .env file in the project root
		env := filepath.Join(getProjectRoot(), ".env")
		if err := godotenv.Load(env); err != nil {
			log.Printf("Error loading .env file: %v", err)
		}

		var psqlInfo string

		if strings.HasSuffix(os.Args[0], ".test") {
			host := os.Getenv("TEST_DB_HOST")
			port := os.Getenv("TEST_DB_PORT")
			user := os.Getenv("TEST_DB_USER")
			password := os.Getenv("TEST_DB_PASSWORD")
			dbname := os.Getenv("TEST_DB_NAME")
			sslmode := os.Getenv("TEST_SSL_MODE")

			psqlInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				host, port, user, password, dbname, sslmode)

		} else {
			host := os.Getenv("DB_HOST")
			port := os.Getenv("DB_PORT")
			user := os.Getenv("DB_USER")
			password := os.Getenv("DB_PASSWORD")
			dbname := os.Getenv("DB_NAME")
			sslmode := os.Getenv("SSL_MODE")

			psqlInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				host, port, user, password, dbname, sslmode)
		}

		var err error
		for i := 0; i < 10; i++ {
			db, err = sql.Open("postgres", psqlInfo)
			if err == nil && db.Ping() == nil {
				log.Printf("Successfully connected to the database!")
				return
			}
			log.Printf("Waiting for database to be ready... (%d/10)", i+1)
			time.Sleep(1 * time.Second)
		}

		log.Fatalf("Could not connect to the database: %v", err)

	})
	return db
}

func getProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}
