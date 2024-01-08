package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	//if err := runMigrations(); err != nil {
	//	log.Fatal(err)
	//}

	// Вернуть nil, если соединение успешно
	return nil
}

//func runMigrations() error {
//	dir := "file://db/migrations"
//	driver, err := mongodb.WithInstance(Client, &mongodb.Config{})
//	if err != nil {
//		return err
//	}
//
//	m, err := migrate.NewWithDatabaseInstance(
//		fmt.Sprintf("%s", dir),
//		"mongodb", driver,
//	)
//	if err != nil {
//		return err
//	}
//
//	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
//		return err
//	}
//
//	fmt.Println("Migrations applied successfully!")
//	return nil
//}
