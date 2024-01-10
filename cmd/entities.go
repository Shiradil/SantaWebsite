package cmd

import "go.mongodb.org/mongo-driver/bson/primitive"

type Volunteer struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Surname  string             `json:"lastName" bson:"lastName"`
	Email    string             `json:"email" bson:"email"`
	Phone    string             `json:"phone" bson:"phone"`
	Password string             `json:"password" bson:"password"`
	Child    *Child             `json:"child,omitempty" bson:"child,omitempty"`
}

type Child struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Surname   string             `json:"surname" bson:"surname"`
	Email     string             `json:"email" bson:"email"`
	Phone     string             `json:"phone" bson:"phone"`
	Password  string             `json:"password" bson:"password"`
	Wish      string             `json:"wish" bson:"wish"`
	Volunteer *Volunteer         `json:"volunteer,omitempty" bson:"volunteer,omitempty"`
}

type WishesData struct {
    ChildId primitive.ObjectID `json:"childId"`
    Wish    string             `json:"wish"`
}

