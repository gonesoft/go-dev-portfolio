package users

import (
	"database/sql"
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

func TestUpdateUser(t *testing.T) {
	testDB := db.Connect()

	var id int
	err := testDB.QueryRow(`INSERT INTO users (name, email) 
		VALUES ($1, $2) RETURNING id`, "Old Nme", "old@example.com").Scan(&id)
	assert.NoError(t, err, "Failed to insert user")

	err = UpdateUserFromDB(testDB, id, "New Name", "new@example.com")
	assert.NoError(t, err, "Failed to update user")

	var udatedName, updatedEmail string
	err = testDB.QueryRow(`SELECT name, email FROM users WHERE id = $1`, id).Scan(&udatedName, &updatedEmail)
	assert.NoError(t, err, "Failed to fetch updated user")
	assert.Equal(t, "New Name", udatedName, "User name was not updated correctly")
	assert.Equal(t, "new@example.com", updatedEmail, "User email was not updated correctly")
}
func TestDeleteUser(t *testing.T) {
	testDB := db.Connect()

	var id int
	err := testDB.QueryRow(`INSERT INTO users (name, email) 
		VALUES ($1, $2) RETURNING id`, "Delete Me", "delete@example.com").Scan(&id)
	assert.NoError(t, err, "Failed to insert user for deletion")

	err = DeleteUserFromDB(testDB, id)
	assert.NoError(t, err, "Failed to delete user")

	var deletedAt sql.NullTime
	err = testDB.QueryRow(`SELECT deleted_at FROM users WHERE id = $1`, id).Scan(&deletedAt)
	assert.NoError(t, err, "Failed to fetch deleted user")
	assert.True(t, deletedAt.Valid, "User should be soft-deleted")
}

func TestGetUserByID(t *testing.T) {
	testDB := db.Connect()

	var id int
	err := testDB.QueryRow(`INSERT INTO users (name, email) 
		VALUES ($1, $2) RETURNING id`, "New User", "newemail@example.com").Scan(&id)
	assert.NoError(t, err, "Failed to insert user for fetching")
	user, err := GetUserByIDFromDB(testDB, id)
	assert.NoError(t, err, "Failed to fetch user by ID")
	assert.Equal(t, "New User", user.Name, "User name does not match")

}

func TestCreateUserInDB(t *testing.T) {
	testDB := db.Connect()

	var user User
	user.Name = "Test User"
	user.Email = "testuser@example.com"
	err := CreateUserInDB(testDB, &user)
	assert.NoError(t, err, "Failed to insert user")
	assert.Greater(t, user.ID, 0, "User ID should be greater than 0")
}
