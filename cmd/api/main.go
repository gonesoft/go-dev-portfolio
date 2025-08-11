package main

import (
	"fmt"
	"log"
	"net/http"

	"gonesoft/go-dev-portfolio/internal/users"
)

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			users.GetUsers(w, r)
		} else if r.Method == http.MethodPost {
			users.CreateUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			users.GetUserByID(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server running on http://localhost:8083")
	log.Fatal(http.ListenAndServe(":8083", nil))

}
