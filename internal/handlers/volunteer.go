package handlers

import (
	"SantaWeb/internal/db"
	"SantaWeb/structs"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func VolLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		phone := r.FormValue("phone")
		password := r.FormValue("password")
		incMsg := "Wrong password or phone"

		collection := db.Client.Database("SantaWeb").Collection("volunteers")
		var volunteer structs.Volunteer
		err := collection.FindOne(context.Background(), bson.M{"phone": phone}).Decode(&volunteer)
		if err != nil {
			RenderTemplate(w, "vollogin.html", incMsg)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(volunteer.Password), []byte(password))
		if err != nil {
			RenderTemplate(w, "vollogin.html", incMsg)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/vol/%s", volunteer.ID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		RenderTemplate(w, "vollogin.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func VolRegHandler(w http.ResponseWriter, r *http.Request) {
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

		volunteer := structs.Volunteer{
			Name:     firstName,
			Surname:  lastName,
			Email:    email,
			Phone:    phone,
			Password: string(hashedPassword),
			Child:    &structs.Child{},
		}

		collection := db.Client.Database("SantaWeb").Collection("volunteers")

		result, err := collection.InsertOne(context.Background(), volunteer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		insertedID := result.InsertedID.(primitive.ObjectID)
		http.Redirect(w, r, fmt.Sprintf("/vol/%s", insertedID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		RenderTemplate(w, "volreg.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func VolunteerPersonalPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	volunteerID := vars["id"]

	var volunteer structs.Volunteer
	collection := db.Client.Database("SantaWeb").Collection("volunteers")
	objID, _ := primitive.ObjectIDFromHex(volunteerID)

	err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&volunteer)
	if err != nil {
		http.Error(w, "Volunteer not found", http.StatusNotFound)
		return
	}

	page := r.URL.Query().Get("page")
	var p int
	if page == "" {
		p = 0
	} else {
		p, _ = strconv.Atoi(page)
	}

	children, _ := GetChildren(p)

	data := struct {
		Volunteer structs.Volunteer
		Children  []structs.Child
	}{
		Volunteer: volunteer,
		Children:  children,
	}

	RenderTemplate(w, "vol.html", data)
}

func GetChildren(page int) ([]structs.Child, error) {
	limit := 10
	offset := 0

	if page > 1 {
		offset = (page - 1) * limit
	}

	collection := db.Client.Database("SantaWeb").Collection("children")

	findOptions := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := collection.Find(context.Background(), bson.D{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error finding children: %v", err)
	}
	defer cursor.Close(context.Background())

	var children []structs.Child
	for cursor.Next(context.Background()) {
		var child structs.Child
		if err := cursor.Decode(&child); err != nil {
			return nil, fmt.Errorf("error decoding children: %v", err)
		}
		children = append(children, child)
	}

	return children, nil
}
