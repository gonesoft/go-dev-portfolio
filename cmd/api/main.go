package main

import (
	"fmt"
	"log"
	"net/http"

	"gonesoft/go-dev-portfolio/internal/users"
)

func main() {

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			users.GetUsers(w, r)
		case http.MethodPost:
			users.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			users.UpdateUser(w, r)
		case http.MethodDelete:
			users.DeleteUser(w, r)
		case http.MethodGet:
			users.GetUserByID(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server running on http://localhost:8083")
	log.Fatal(http.ListenAndServe(":8083", nil))

}
