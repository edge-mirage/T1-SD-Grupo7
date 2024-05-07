package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Last_name string             `bson:"last_name" json:"last_name"`
	Rut       string             `bson:"rut" json:"rut"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
}

type RegisterRequest struct {
	Name      string `bson:"name" json:"name"`
	Last_name string `bson:"last_name" json:"last_name"`
	Rut       string `bson:"rut" json:"rut"`
	Email     string `bson:"email" json:"email"`
	Password  string `bson:"password" json:"password"`
}

type LoginRequest struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}
