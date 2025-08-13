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

	database := db.Connect()
	if database == nil {
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
	database := db.Connect()
	if database == nil {
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
	database := db.Connect()
	if database == nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer database.Close()

	var u User
	database.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)

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
	database := db.Connect()
	if database == nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer database.Close()

	database.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&user.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUsers handles GET /users request :)
func GetUsers(w http.ResponseWriter, r *http.Request) {

	//Pagination
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	//searching
	searchTerm := r.URL.Query().Get("search")
	if searchTerm != "" {
		// Sanitize search term to prevent SQL injection
		searchTerm = strings.TrimSpace(searchTerm)
		searchTerm = "%" + strings.ToLower(searchTerm) + "%"
	} else {
		searchTerm = "%"
	}

	// Sorting. If no sort parameter is provided, default to sorting by name
	sortBy := r.URL.Query().Get("sort")
	if sortBy == "" && sortBy != "email" && sortBy != "created_at" {
		sortBy = "name"
	}

	//ordering
	order := strings.ToUpper(r.URL.Query().Get("order"))
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	// Connect to the database
	database := db.Connect()
	if database == nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to connect to the database")
		return
	}
	defer database.Close()

	usersList, total, err := GetUsersFromDB(database, searchTerm, limit, offset, sortBy, order)
	if err != nil {
		httphelper.Error(w, http.StatusInternalServerError, "Failed to fetch users: "+err.Error())
		return
	}

	totalPages := (total + limit - 1) / limit // Calculate total pages
	httphelper.JSON(w, http.StatusOK, map[string]interface{}{
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
		"data":        usersList,
	})
}
