package main

import (
	"SantaWeb/internal/db"
	"SantaWeb/internal/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	logFile, _ := os.OpenFile("cmd/logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()
	log.SetOutput(logFile)

	handlers.InitLogger(log)

	err := db.DbConnection()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	router := mux.NewRouter()
	handlers.SetupRoutes(router)

	port := ":8080"

	log.Infof("Starting server on http://localhost:%s/\n", port)
	fmt.Printf("Starting server on http://localhost%s/\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
