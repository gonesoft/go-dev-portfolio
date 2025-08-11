package main

import (
	"fmt"
	"log"
	"net/http"

	"gonesoft/go-dev-portfolio/internal/users"
)

func main() {
	http.HandleFunc("/users", users.GetUsers)

	fmt.Println("Starting server on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))

}
