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

const pageSize = 10

type PaginationData struct {
	CurrentPage int
	PrevPage    int
	NextPage    int
	TotalPages  int
	Pages       []int
}

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

	volunteer, err := GetVolunteerByID(volunteerID)
	if err != nil {
		http.Error(w, "Volunteer not found", http.StatusNotFound)
		return
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	children, totalCount, err := GetChildren(page)
	if err != nil {
		http.Error(w, "Error retrieving children data", http.StatusInternalServerError)
		return
	}

	pagination := CalculatePagination(page, totalCount)

	data := struct {
		Volunteer  structs.Volunteer
		Children   []structs.Child
		Pagination PaginationData
	}{
		Volunteer:  volunteer,
		Children:   children,
		Pagination: pagination,
	}

	RenderTemplate(w, "vol.html", data)
}

// http://localhost:8080/vol/65bdf00b869485d29a4c66e0?page=2
// http://localhost:8080/vol/ObjectID%28%2265bdf00b869485d29a4c66e0%22%29?page=1
func GetVolunteerByID(volunteerID string) (structs.Volunteer, error) {
	var volunteer structs.Volunteer

	objID, err := primitive.ObjectIDFromHex(volunteerID)
	if err != nil {
		return volunteer, fmt.Errorf("invalid volunteer ID: %v", err)
	}

	collection := db.Client.Database("SantaWeb").Collection("volunteers")
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&volunteer)
	if err != nil {
		return volunteer, fmt.Errorf("error finding volunteer: %v", err)
	}

	return volunteer, nil
}

func CalculatePagination(page, totalCount int) PaginationData {
	totalPages := (totalCount + pageSize - 1) / pageSize
	prevPage := page - 1
	nextPage := page + 1

	return PaginationData{
		CurrentPage: page,
		PrevPage:    prevPage,
		NextPage:    nextPage,
		TotalPages:  totalPages,
	}
}

func GetChildren(page int) ([]structs.Child, int, error) {
	limit := 10
	offset := (page - 1) * limit

	collection := db.Client.Database("SantaWeb").Collection("children")
	ctx := context.Background()

	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)))
	if err != nil {
		return nil, 0, fmt.Errorf("error finding children: %v", err)
	}
	defer cursor.Close(ctx)

	var children []structs.Child
	for cursor.Next(ctx) {
		var child structs.Child
		if err := cursor.Decode(&child); err != nil {
			return nil, 0, fmt.Errorf("error decoding children: %v", err)
		}
		children = append(children, child)
	}

	totalCount, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count of children: %v", err)
	}

	return children, int(totalCount), nil
}
