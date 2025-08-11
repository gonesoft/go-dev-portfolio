package users

import (
	"encoding/json"
	"net/http"

	"gonesoft/go-dev-portfolio/internal/db"
)

// GetUsers handles GET /users request :)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Connect to the database
	database, err := db.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	// Query the database for users
	rows, err := database.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, "Failed to query users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var allUsers []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}
		allUsers = append(allUsers, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allUsers)
}
