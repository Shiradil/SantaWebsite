package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os/user"
)

var Client *mongo.Client

func DbConnection() error {
	clientOptions := options.Client().ApplyURI("mongodb+srv://shirbaev04:bauka@cluster0.xttjkma.mongodb.net/")

	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	err = Client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	fmt.Println("Connected to MongoDB!")
	id, err := primitive.ObjectIDFromHex("65980c087ef3cf7bf3fc6870")

	collection := Client.Database("SantaWeb").Collection("volunteers")
	var result user.User
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("Документ не найден")
		} else {
			log.Fatal(err)
		}
	}

	fmt.Println(result.Username)
	// Вернуть nil, если соединение успешно
	return nil
}
