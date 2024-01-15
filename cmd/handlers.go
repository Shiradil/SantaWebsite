package cmd

import (
	"SantaWeb/db"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

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
		phone := r.FormValue("phone")
		password := r.FormValue("password")

		collection := db.Client.Database("SantaWeb").Collection("children")
		var child Child
		err := collection.FindOne(context.Background(), bson.M{"phone": phone}).Decode(&child)
		if err != nil {
			errorResponse := ErrorResponse{Status: "400", Message: "Некорректное JSON-сообщение"}
			sendJSONResponse(w, errorResponse, http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(child.Password), []byte(password))
		if err != nil {
			errorResponse := ErrorResponse{Status: "400", Message: "Некорректное JSON-сообщение"}
			sendJSONResponse(w, errorResponse, http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/chil/%s", child.ID.Hex()), http.StatusSeeOther)
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

		result, err := collection.InsertOne(context.Background(), child)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		insertedID := result.InsertedID.(primitive.ObjectID)
		http.Redirect(w, r, fmt.Sprintf("/chil/%s", insertedID.Hex()), http.StatusSeeOther)
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

func updateWishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		childIdString := r.FormValue("childId")
		newWish := r.FormValue("wish")

		childId, err := primitive.ObjectIDFromHex(childIdString)
		if err != nil {
			http.Error(w, "Invalid Child ID", http.StatusBadRequest)
			return
		}

		collection := db.Client.Database("SantaWeb").Collection("Children")

		filter := bson.M{"_id": childId}
		update := bson.M{"$set": bson.M{"wish": newWish}}
		_, err = collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/path-after-updating-wish", http.StatusSeeOther)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
