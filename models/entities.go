package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Volunteer struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Surname     string             `json:"lastName" bson:"lastName"`
	Email       string             `json:"email" bson:"email"`
	Phone       string             `json:"phone" bson:"phone"`
	Password    string             `json:"password" bson:"password"`
	Child       *Child             `json:"child,omitempty" bson:"child,omitempty"`
	ConfirmCode string                `json:"confirmCode,omitempty" bson:"confirmCode,omitempty"`
	IsConfirmed bool               `json:"isConfirmed" bson:"isConfirmed"`
}

type Child struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Surname   string             `json:"surname" bson:"surname"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	Password  string             `json:"password" bson:"password"`
	Wish      *Wish              `json:"wish,omitempty" bson:"wish,omitempty"`
	Volunteer *Volunteer         `json:"volunteer,omitempty" bson:"volunteer,omitempty"`
}

type Wish struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Wishes  string             `json:"wishes" bson:"wishes"`
	ChildID primitive.ObjectID `json:"childID,omitempty" bson:"childID,omitempty"`
}
