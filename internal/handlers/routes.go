package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("ui/static"))))

	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/vol/{id}", VolunteerPersonalPageHandler)
	router.HandleFunc("/chil/{id}", ChildPersonalPageHandler)
	router.HandleFunc("/vollogin", VolLoginHandler)
	router.HandleFunc("/volreg", VolRegHandler)
	router.HandleFunc("/chilog", ChiLogHandler)
	router.HandleFunc("/chireg", ChiRegHandler)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	})
	router.HandleFunc("/update-wishes", UpdateWishesHandler).Methods("POST")
}
