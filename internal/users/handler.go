package users

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gonesoft/go-dev-portfolio/internal/db"
	httphelper "gonesoft/go-dev-portfolio/internal/http"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		httphelper.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		httphelper.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if user.Name == "" || user.Email == "" {
		httphelper.Error(w, http.StatusBadRequest, "Name and email are required")
		return
	}

	database, err := db.Connect()
	if err != nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer database.Close()

	result, err := database.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", user.Name, user.Email, id)
	if err != nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to update user")
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		httphelper.Error(w, http.StatusNotFound, "User not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		httphelper.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	database, err := db.Connect()
	if err != nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer database.Close()
	result, err := database.Exec("UPDATE users SET deleted_at = NOW() WHERE id = $1", id)
	if err != nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		httphelper.Error(w, http.StatusNotFound, "User not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		httphelper.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	database, err := db.Connect()
	if err != nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer database.Close()

	var u User
	err = database.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		httphelper.Error(w, http.StatusInternalServerError, "User not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		httphelper.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if user.Name == "" || user.Email == "" {
		httphelper.Error(w, http.StatusBadRequest, "Name and email are required")
		return
	}

	// Connect to the database
	database, err := db.Connect()
	if err != nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer database.Close()

	err = database.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUsers handles GET /users request :)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Connect to the database
	database, err := db.Connect()
	if err != nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer database.Close()

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Query the database for users
	rows, err := database.Query(`SELECT id, name, email FROM users WHERE deleted_at 
	IS NULL LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to query users: "+err.Error())
		return
	}
	defer rows.Close()

	var allUsers []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			httphelper.Error(w, http.StatusInternalServerError, "Failed to scan user")
			return
		}
		allUsers = append(allUsers, user)
	}

	if err := rows.Err(); err != nil {
		httphelper.Error(w, http.StatusInternalServerError, "Error processing users")
		return
	}

	httphelper.JSON(w, http.StatusOK, map[string]interface{}{
		"page":  page,
		"limit": limit,
		"data":  allUsers,
	})
}
