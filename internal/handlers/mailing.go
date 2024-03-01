package handlers

import (
	"SantaWeb/internal/db"
	"SantaWeb/services"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"net/http"
)

type requestBody struct {
	Text string
}

func MailingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var requestBody requestBody
	if err := json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	text := requestBody.Text
	collection := db.Client.Database("SantaWeb").Collection("volunteers")
	filter := bson.M{}

	var emails []string

	cursor, err := collection.Find(context.Background(), filter, options.Find().SetProjection(bson.D{{"email", 1}}))
	if err != nil {
		log.Error(err.Error())
		return
	}

	for cursor.Next(context.Background()) {
		var email struct {
			Email string `bson:"email"`
		}
		if err := cursor.Decode(&email); err != nil {
			fmt.Println("Error decoding document:", err)
			return
		}
		emails = append(emails, email.Email)
	}
	err = cursor.Close(context.Background())
	if err != nil {
		log.Error(err.Error())
		return
	}

	go func() {
		for _, email := range emails {
			services.SendMail(email, text)
		}
	}()
}
