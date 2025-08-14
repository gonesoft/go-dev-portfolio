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

	// 1) Create a user
	resp, err := http.Post(
		ts.URL+"/users",
		"application/json",
		strings.NewReader(`{"name":"Test Creation","email":"testcreation@example.com"}`),
	)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	var createdUser User
	_ = json.NewDecoder(resp.Body).Decode(&createdUser)
	resp.Body.Close()

	// 2) Fetch the user by ID
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(createdUser.ID))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// 3) Update the user (correct URL with slash)
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

	// 4) Delete the user
	req, _ = http.NewRequest(http.MethodDelete, ts.URL+"/users/"+strconv.Itoa(createdUser.ID), nil)
	delResp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, delResp.StatusCode)
	delResp.Body.Close()

	// 5) Verify user is deleted -> 404
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(createdUser.ID))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()
}
