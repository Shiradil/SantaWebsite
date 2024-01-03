package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RunServer() {
	router := mux.NewRouter()

	router.HandleFunc("/", HomeHandler).Methods("GET")
	router.HandleFunc("/letter", LetterHandler).Methods("POST")

	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Добро пожаловать на страницу написания письма Санте!")
}

func LetterHandler(w http.ResponseWriter, r *http.Request) {
	// тут напишем добавление писем в базу данных

	fmt.Fprintln(w, "Спасибо за ваше письмо! Санта обязательно его увидит.")
}
