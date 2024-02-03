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
		w.WriteHeader(http.StatusNotFound)
		ErrorHandler(w, r, http.StatusNotFound, "Volunteer not found")
		return
	}

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	var sortDirection int
	sortParam := r.URL.Query().Get("sort")
	if sortParam == "asc" {
		sortDirection = 1
	} else {
		sortDirection = -1
	}

	children, totalCount, err := GetChildren(page, sortDirection)
	if err == fmt.Errorf("page does not exist") {
		w.WriteHeader(http.StatusNotFound)
		ErrorHandler(w, r, http.StatusNotFound, "Volunteer not found")
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	pagination := CalculatePagination(page, totalCount)

	data := struct {
		Volunteer  structs.Volunteer
		Children   []structs.Child
		Pagination PaginationData
		Sorting    string
	}{
		Volunteer:  volunteer,
		Children:   children,
		Pagination: pagination,
		Sorting:    sortParam,
	}

	RenderTemplate(w, "vol.html", data)
}

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

func GetChildren(page int, sortDirection int) ([]structs.Child, int, error) {
	limit := 10
	offset := (page - 1) * limit

	collection := db.Client.Database("SantaWeb").Collection("children")
	ctx := context.Background()

	sort := bson.D{{Key: "name", Value: sortDirection}}

	cursor, err := collection.Find(ctx, bson.D{}, options.Find().SetLimit(int64(limit)).SetSkip(int64(offset)).SetSort(sort))
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

	totalCount, err := collection.EstimatedDocumentCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count of children: %v", err)
	}

	totalPages := (int(totalCount) + limit - 1) / limit
	if page > totalPages {
		return nil, totalPages, fmt.Errorf("page does not exist")
	}

	return children, totalPages, nil
}
