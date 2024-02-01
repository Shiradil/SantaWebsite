package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes() error {
	// Установка опций клиента
	clientOptions := options.Client().ApplyURI("mongodb+srv://shirbaev04:bauka@cluster0.xttjkma.mongodb.net/")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	// Создание индексов для коллекции "volunteers"
	volunteersCollection := client.Database("SantaWeb").Collection("volunteers")
	_, err = volunteersCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    map[string]interface{}{"phone": 1},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return err
	}

	// Создание индексов для коллекции "children"
	childrenCollection := client.Database("SantaWeb").Collection("children")
	_, err = childrenCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    map[string]interface{}{"phone": 1},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
