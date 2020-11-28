package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser = "users"
)

type User struct {
	Username  string             `json:"username" bson:"username" binding:"required"`
	Password  string             `json:"password" bson:"password" binding:"required"`
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	Email     string             `json:"email" bson:"email" binding:"required,email"`
}

type UserOption struct {
	Password string `json:"password" bson:"password" binding:"omitempty"`
	OldPassword string `json:"old_password" bson:"-" binding:"required_with=Password"`
	Email    string `json:"email" bson:"email" binding:"omitempty,email"`
}

var UserProjection = bson.D{
	{"username", 1},
	{"_id", 0},
	{"id", "$_id"},
	{"created_at", 1},
	{"email", 1},
}
