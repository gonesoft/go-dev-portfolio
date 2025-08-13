package users

import (
	"gonesoft/go-dev-portfolio/internal/db"
	"log"
	"os"
	"testing"

	//"github.com/go-playground/assert/v2"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	testDB := db.Connect()
	if testDB == nil {
		log.Fatal("Failed to connect to the test database")
	}
	defer testDB.Close()

	_, _ = testDB.Exec("DELETE FROM users")

	os.Exit(m.Run())
}

func TestCreateUserAndFetch(t *testing.T) {
	testDB := db.Connect()
	_, err := testDB.Exec(`INSERT INTO users (name, email) VALUES ($1, $2)`, "Test User", "test@example4.com")
	assert.NoError(t, err, "Failed to insert user")

	usersList, total, err := GetUsersFromDB(testDB, "Test User", 10, 0, "name", "ASC")
	assert.NoError(t, err, "Failed to fetch users")
	assert.Equal(t, 1, total, "Expected 1 user to be returned")
	assert.Equal(t, "Test User", usersList[0].Name, "User name does not match")
}
