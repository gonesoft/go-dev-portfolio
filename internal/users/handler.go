package users

import (
	"database/sql"
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
	//defer database.Close()

	err = UpdateUserFromDB(database, id, &user)
	if err != nil {
		if err == sql.ErrNoRows {
			httphelper.Error(w, http.StatusNotFound, "User not found")
			return
		}
		if strings.Contains(err.Error(), "exists") {
			httphelper.Error(w, http.StatusConflict, "Email already exists")
			return
		}
		httphelper.Error(w, http.StatusInternalServerError, "Failed to update user")
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
	//defer database.Close()

	err = DeleteUserFromDB(database, id)
	if err != nil {
		if err == sql.ErrNoRows {
			httphelper.Error(w, http.StatusNotFound, "User not found")
			return
		}
		httphelper.Error(w, http.StatusInternalServerError, "Failed to delete user")
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
	//defer database.Close()

	var user User
	user, err = GetUserByIDFromDB(database, id)
	if err != nil {
		httphelper.Error(w, http.StatusNotFound, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

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
	//defer database.Close()

	err := CreateUserInDB(database, &user)
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
	//defer database.Close()

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

func GetUsersNoPaging(w http.ResponseWriter, r *http.Request) {
	db := db.Connect()

	// Parse query params
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	sortBy := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	search := r.URL.Query().Get("search")

	opts := ListOptions{
		Search: search,
		Limit:  limit,
		Offset: offset,
		SortBy: sortBy,
		Order:  order,
	}

	list, total, err := ListUsers(db, opts)
	if err != nil {
		switch err {
		case ErrInvalidSort:
			http.Error(w, "invalid sort: allowed id,name,email,created_at", http.StatusBadRequest)
			return
		case ErrInvalidOrder:
			http.Error(w, "invalid order: allowed ASC,DESC", http.StatusBadRequest)
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	// compute total pages if client passed a limit (avoid div by zero)
	totalPages := 0
	if opts.Limit > 0 {
		totalPages = (total + opts.Limit - 1) / opts.Limit
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"limit":       opts.Limit,
		"offset":      opts.Offset,
		"total":       total,
		"total_pages": totalPages,
		"sort":        strings.ToLower(opts.SortBy),
		"order":       strings.ToUpper(opts.Order),
		"users":       list,
	})
}
