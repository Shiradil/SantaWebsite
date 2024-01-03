package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func DbConnection() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	//collection := client.Database("SantaWeb").Collection("wishes")
	// тут крч будут функции для работы с монгодб
}
