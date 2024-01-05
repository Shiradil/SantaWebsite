package cmd

import (
	"SantaWeb/db"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Volunteer struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Surname  string             `bson:"lastName"`
	Email    string             `bson:"email"`
	Phone    string             `bson:"phone"`
	Password string             `bson:"password"`
	Child    *Child             `bson:"child,omitempty"`
}

type Child struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Surname   string             `bson:"surname"`
	Email     string             `bson:"name"`
	Phone     string             `bson:"name"`
	Password  string             `bson:"name"`
	Wish      string             `bson:"wish"`
	Volunteer *Volunteer         `bson:"volunteer,omitempty"`
}

func RunServer() {
	// подключение к монгодб
	db.DbConnection()

	router := mux.NewRouter()

	// templates connecting
	router.HandleFunc("/", PageHandler("home.html")).Methods("GET")
	// children
	router.HandleFunc("/chilog", PageHandler("chilog.html")).Methods("GET")
	router.HandleFunc("/chireg", PageHandler("chireg.html")).Methods("GET")
	router.HandleFunc("/chil", PageHandler("chil.html")).Methods("GET")
	// volunteers
	router.HandleFunc("/vollogin", PageHandler("vollogin.html")).Methods("GET")
	router.HandleFunc("/volreg", PageHandler("volreg.html")).Methods("GET")
	router.HandleFunc("/vol", PageHandler("vol.html")).Methods("GET")

	// css and js connecting
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	router.HandleFunc("/letter", LetterHandler).Methods("POST")
	router.HandleFunc("/submit-volunteer-registration", RegisterVolunteerHandler).Methods("POST")
	router.HandleFunc("/submit-child-registration", RegisterChildHandler).Methods("POST")

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

func RegisterVolunteerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	volunteer := Volunteer{
		Name:     firstName,
		Surname:  lastName,
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
		Child:    nil,
	}

	collection := db.Client.Database("SantaWeb").Collection("volunteers")
	_, err = collection.InsertOne(context.Background(), volunteer)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Registration Successful!"))
}

func RegisterChildHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	child := Child{
		Name:      firstName,
		Surname:   lastName,
		Email:     email,
		Phone:     phone,
		Password:  string(hashedPassword),
		Wish:      "",
		Volunteer: nil,
	}

	collection := db.Client.Database("SantaWeb").Collection("children")
	_, err = collection.InsertOne(context.Background(), child)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Registration Successful!"))
}
