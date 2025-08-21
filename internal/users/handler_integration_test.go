package users

import (
	"encoding/json"
	"gonesoft/go-dev-portfolio/internal/db"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserLifeCycle(t *testing.T) {
	testDB := db.Connect()
	_, _ = testDB.Exec("DELETE FROM users")

	// Router wiring:
	// - /users       -> GET (list), POST (create)
	// - /users/{id}  -> GET (single), PUT (update), DELETE (soft delete)
	mux := http.NewServeMux()
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetUsers(w, r)
		case http.MethodPost:
			CreateUser(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			GetUserByID(w, r)
		case http.MethodPut:
			UpdateUser(w, r)
		case http.MethodDelete:
			DeleteUser(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	var createdUser User

	// 1) Create a user
	t.Run("Create User", func(t *testing.T) {
		resp, err := http.Post(
			ts.URL+"/users",
			"application/json",
			strings.NewReader(`{"name":"Test Creation","email":"testcreation@example.com"}`),
		)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		_ = json.NewDecoder(resp.Body).Decode(&createdUser)
		resp.Body.Close()
	})

	// 2) Fetch the user by ID
	t.Run("Fetch User", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/users/" + strconv.Itoa(createdUser.ID))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	})

	// 3) Update the user (correct URL with slash)
	t.Run("Update User", func(t *testing.T) {
		req, _ := http.NewRequest(
			http.MethodPut,
			ts.URL+"/users/"+strconv.Itoa(createdUser.ID),
			strings.NewReader(`{"name":"Updated","email":"updated@example.com"}`),
		)
		req.Header.Set("Content-Type", "application/json")
		updatedResp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		// Expect 204 No Content if your handler follows that convention
		assert.Equal(t, http.StatusNoContent, updatedResp.StatusCode)
		updatedResp.Body.Close()
	})

	// 4) Delete the user
	t.Run("Delete User", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/users/"+strconv.Itoa(createdUser.ID), nil)
		delResp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, delResp.StatusCode)
		delResp.Body.Close()
	})

	// 5) Verify user is deleted -> 404
	t.Run("Verify User Deletion", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/users/" + strconv.Itoa(createdUser.ID))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		resp.Body.Close()
	})

	// 6) Test email duplication
	t.Run("Test Email Duplication expected 409", func(t *testing.T) {
		_, _ = http.Post(ts.URL+"/users", "application/json", strings.NewReader(
			`{"name":"User1","email":"user1@example.com"}`))
		_, _ = http.Post(ts.URL+"/users", "application/json", strings.NewReader(
			`{"name":"User2","email":"user2@example.com"}`))

		req, _ := http.NewRequest(http.MethodPut, ts.URL+"/users/2",
			strings.NewReader(`{"name":"User2 Updated","email":"user1@example.com"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	// 7) Test update non-existent user
	t.Run("Update non-existent user returns 404", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, ts.URL+"/users/9999",
			strings.NewReader(`{"name":"Doesn't exist","email":"nope@example.com"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := http.DefaultClient.Do(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	// 8) Test delete non-existent user
	t.Run("Delete non-existent user returns 404", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, ts.URL+"/users/9999", nil)
		resp, _ := http.DefaultClient.Do(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	//9) List Users with Paging
	t.Run("List Users with Paging", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/users?limit=5&offset=0")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		var users []User
		_ = json.NewDecoder(resp.Body).Decode(&users)
		assert.LessOrEqual(t, len(users), 5, "Expected no more than 5 users in the response")
	})
}
