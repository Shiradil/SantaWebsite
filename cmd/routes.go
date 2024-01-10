package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
)

func setupRoutes(router *mux.Router) {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/vol/{id}", volunteerPersonalPageHandler).Methods("GET")
	router.HandleFunc("/chil/{id}", childPersonalPageHandler).Methods("GET")
	router.HandleFunc("/vollogin", volLoginHandler).Methods("GET", "POST")
	router.HandleFunc("/volreg", volRegHandler).Methods("GET", "POST")
	router.HandleFunc("/chilog", chiLogHandler).Methods("GET", "POST")
	router.HandleFunc("/chireg", chiRegHandler).Methods("GET", "POST")
	router.HandleFunc("/update-wish", updateWishHandler).Methods("PUT")
}
