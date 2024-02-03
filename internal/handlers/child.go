package handlers

import (
	"SantaWeb/internal/db"
	"SantaWeb/structs"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func ChiLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		phone := r.FormValue("phone")
		password := r.FormValue("password")
		incMsg := "Wrong password or phone"

		collection := db.Client.Database("SantaWeb").Collection("children")
		var child structs.Child
		err := collection.FindOne(context.Background(), bson.M{"phone": phone}).Decode(&child)
		if err != nil {
			RenderTemplate(w, "chilog.html", incMsg)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(child.Password), []byte(password))
		if err != nil {
			RenderTemplate(w, "chilog.html", incMsg)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/chil/%s", child.ID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		RenderTemplate(w, "chilog.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func ChiRegHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		password := r.FormValue("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		child := structs.Child{
			Name:      firstName,
			Surname:   lastName,
			Email:     email,
			Phone:     phone,
			Password:  string(hashedPassword),
			Wish:      &structs.Wish{},
			Volunteer: &structs.Volunteer{},
		}

		collection := db.Client.Database("SantaWeb").Collection("children")
		collectionWishes := db.Client.Database("SantaWeb").Collection("wishes")

		result, err := collection.InsertOne(context.Background(), child)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			return
		}

		insertedID := result.InsertedID.(primitive.ObjectID)

		wishes := structs.Wish{
			Wishes: "",
			Child:  &child,
		}
		_, err = collectionWishes.InsertOne(context.Background(), wishes)
		if err != nil {
			http.Error(w, "Error creating wish", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/chil/%s", insertedID.Hex()), http.StatusSeeOther)
	} else if r.Method == "GET" {
		RenderTemplate(w, "chireg.html", nil)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func ChildPersonalPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	childID := vars["id"]

	var child structs.Child
	collection := db.Client.Database("SantaWeb").Collection("children")
	objID, _ := primitive.ObjectIDFromHex(childID)

	err := collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&child)
	if err != nil {
		http.Error(w, "Child not found", http.StatusNotFound)
		return
	}

	RenderTemplate(w, "chil.html", child)
}

func UpdateWishesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	childID, err := GetChildIDFromSession(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	wishes := r.FormValue("wishes")

	wishesCollection := db.Client.Database("SantaWeb").Collection("wishes")

	filter := bson.M{"childId": childID}
	update := bson.M{"$set": bson.M{"wishes": wishes}}
	_, err = wishesCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error updating wishes", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/path-to-child-page", http.StatusSeeOther)
}

func GetChildIDFromSession(r *http.Request) (primitive.ObjectID, error) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		return primitive.NilObjectID, err
	}

	childID, ok := session.Values["childID"].(string)
	if !ok {
		return primitive.NilObjectID, errors.New("Child ID not found in session")
	}

	objID, err := primitive.ObjectIDFromHex(childID)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return objID, nil
}
