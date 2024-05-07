package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Token string             `bson:"token" json:"token"`
}

type CreateTokenRequest struct {
	Token string `bson:"token" json:"token"`
}
