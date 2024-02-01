package main

import (
	"SantaWeb/internal/db"
	"SantaWeb/internal/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	err := db.DbConnection()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	router := mux.NewRouter()
	handlers.SetupRoutes(router)

	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
