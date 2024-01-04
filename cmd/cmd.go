package cmd

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RunServer() {
	router := mux.NewRouter()

	// templates
	router.HandleFunc("/", PageHandler("home.html")).Methods("GET")
	// children
	router.HandleFunc("/chilog", PageHandler("chilog.html")).Methods("GET")
	router.HandleFunc("/chireg", PageHandler("chireg.html")).Methods("GET")
	router.HandleFunc("/chil", PageHandler("chil.html")).Methods("GET")
	//volunteers
	router.HandleFunc("/vollogin", PageHandler("vollogin.html")).Methods("GET")
	router.HandleFunc("/volreg", PageHandler("volreg.html")).Methods("GET")
	router.HandleFunc("/vol", PageHandler("vol.html")).Methods("GET")

	//router.HandleFunc("/chilog")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))
	router.HandleFunc("/letter", LetterHandler).Methods("POST")

	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func PageHandler(pageName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("frontend/templates/" + pageName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)
	}
}

func LetterHandler(w http.ResponseWriter, r *http.Request) {
	// тут напишем добавление писем в базу данных

	fmt.Fprintln(w, "Спасибо за ваше письмо! Санта обязательно его увидит.")
}
