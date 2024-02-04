package handlers

import (
	"SantaWeb/internal/db"
	"SantaWeb/structs"
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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
			log.Error(err.Error())
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(child.Password), []byte(password))
		if err != nil {
			RenderTemplate(w, "chilog.html", incMsg)
			log.Error(err.Error())
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
			log.Error(err.Error())
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
			log.Error(err.Error())
			return
		}

		insertedID := result.InsertedID.(primitive.ObjectID)

		wishes := structs.Wish{
			Wishes:  "",
			ChildID: insertedID,
		}
		_, err = collectionWishes.InsertOne(context.Background(), wishes)
		if err != nil {
			http.Error(w, "Error creating wish", http.StatusInternalServerError)
			log.Error(err.Error())
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
		log.Error(err.Error())
		return
	}

	RenderTemplate(w, "chil.html", child)
}
func extractObjectID(input string) (primitive.ObjectID, error) {
	re := regexp.MustCompile(`ObjectID\("([a-fA-F0-9]+)"\)`)
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		return primitive.NilObjectID, fmt.Errorf("no ObjectID found")
	}

	return primitive.ObjectIDFromHex(matches[1])
}

func UpdateWishesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	wishesCollection := db.Client.Database("SantaWeb").Collection("wishes")

	if err := r.ParseForm(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		log.Error(err.Error())
		return
	}

	wishes := r.FormValue("wishes")
	input := r.FormValue("childID")
	childID, err := extractObjectID(input)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error extracting ObjectID: %v", err), http.StatusBadRequest)
		log.Error(err.Error())
		return
	}

	filter := bson.M{"childID": childID}
	update := bson.M{"$set": bson.M{"wishes": wishes}}

	if _, err := wishesCollection.UpdateOne(context.Background(), filter, update); err != nil {
		http.Error(w, "Error updating wishes", http.StatusInternalServerError)
		log.Error(err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/chil/%s", childID.Hex()), http.StatusSeeOther)
}

/*func UpdateWishesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	wishesCollection := db.Client.Database("SantaWeb").Collection("wishes")

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorHandler(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	wishes := r.FormValue("wishes")
	input := r.FormValue("childID")
	childID, err := extractObjectID(input)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Extracted ObjectID:", childID)
	}
	fmt.Println("Child ID:", childID)
	fmt.Println("wishes:", wishes)

	filter := bson.M{"childID": input}
	update := bson.M{"$set": bson.M{"wishes": wishes}}
	_, err = wishesCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, "Error updating wishes")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/chil/%s", childID), http.StatusSeeOther)
}

func extractObjectID(input string) (string, error) {
	re := regexp.MustCompile(`ObjectID\("([a-fA-F0-9]+)"\)`)
	matches := re.FindStringSubmatch(input)

	if len(matches) < 2 {
		return "", fmt.Errorf("no ObjectID found")
	}

	return matches[1], nil
}*/

/*func GetChildIDFromSession(r *http.Request) (primitive.ObjectID, error) {
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
*/
