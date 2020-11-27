package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Username  string             `json:"username" bson:"username" binding:"required"`
	Password  string             `json:"password" bson:"password" binding:"required"`
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	Email     string             `json:"email" bson:"email" binding:"required"`
}
