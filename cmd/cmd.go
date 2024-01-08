package cmd

import (
	"SantaWeb/db"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

type Volunteer struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Surname  string             `json:"lastName" bson:"lastName"`
	Email    string             `json:"email" bson:"email"`
	Phone    string             `json:"phone" bson:"phone"`
	Password string             `json:"password" bson:"password"`
	Child    *Child             `json:"child,omitempty" bson:"child,omitempty"`
}

type Child struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Surname   string             `json:"surname" bson:"surname"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	Password  string             `json:"password" bson:"password"`
	Wish      string             `json:"wish" bson:"wish"`
	Volunteer *Volunteer         `json:"volunteer,omitempty" bson:"volunteer,omitempty"`
}

func RunServer() {
	err := db.DbConnection()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	router := mux.NewRouter()
	setupRoutes(router)

	port := ":8080"
	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func setupRoutes(router *mux.Router) {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))

	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/vol/{id}", volunteerPersonalPageHandler).Methods("GET")
	router.HandleFunc("/chil/{id}", childPersonalPageHandler).Methods("GET")
	router.HandleFunc("/vollogin", volLoginHandler).Methods("GET", "POST")
	router.HandleFunc("/volreg", volRegHandler).Methods("GET", "POST")
	router.HandleFunc("/chilog", chiLogHandler).Methods("GET", "POST")
	router.HandleFunc("/chireg", chiRegHandler).Methods("GET", "POST")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.html", nil)
}

func volLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		phone := r.FormValue("phone")
		password := r.FormValue("password")

		collection := db.Client.Database("SantaWeb").Collection("volunteers")
		var volunteer Volunteer
		err := collection.FindOne(context.Background(), bson.M{"phone": phone}).Decode(&volunteer)
		if err != nil {
			http.Error(w, "Invalid phone or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(volunteer.Password), []byte(password))
		if err != nil {
			http.Error(w, "Invalid phone or password", http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/vol/%s", volunteer.ID.Hex()), http.StatusSeeOther)
	} else {
		renderTemplate(w, "vollogin.html", nil)
	}
}

func volRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
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
		result, err := collection.InsertOne(context.Background(), volunteer)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		insertedID := result.InsertedID.(primitive.ObjectID)
		http.Redirect(w, r, fmt.Sprintf("/vol/%s", insertedID.Hex()), http.StatusSeeOther)
	} else {
		renderTemplate(w, "volreg.html", nil)
	}
}

func chiLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
	} else {
		renderTemplate(w, "chilog.html", nil)
	}
}

func chiRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
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
		http.Redirect(w, r, fmt.Sprintf("/vol/%s", child.ID.Hex()), http.StatusSeeOther)
	} else {
		renderTemplate(w, "chireg.html", nil)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("frontend/templates/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func volunteerPersonalPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	volunteerID := vars["id"]

	var volunteer Volunteer
	collection := db.Client.Database("SantaWeb").Collection("volunteers")
	objID, _ := primitive.ObjectIDFromHex(volunteerID)

	err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&volunteer)
	if err != nil {
		http.Error(w, "Volunteer not found", http.StatusNotFound)
		return
	}

	renderTemplate(w, "vol.html", volunteer)
}

func childPersonalPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	childID := vars["id"]

	var child Child
	collection := db.Client.Database("SantaWeb").Collection("children")
	objID, _ := primitive.ObjectIDFromHex(childID)

	err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&child)
	if err != nil {
		http.Error(w, "Child not found", http.StatusNotFound)
		return
	}

	renderTemplate(w, "chil.html", child)
}
