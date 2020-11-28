package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionClip = "clips"
)

type ClipType int

const (
	Image ClipType = 1
	Text  ClipType = 2
	File  ClipType = 3
)

type Clip struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      ClipType           `bson:"type" json:"type" binding:"required,oneof=1 2 3"`
	Content   string             `bson:"content" json:"content" binding:"required"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	IsDeleted bool               `bson:"is_deleted" json:"is_deleted"`
	DeletedAt time.Time          `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	UserId    primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
}

type ClipOption struct {
	Type    ClipType `bson:"type,omitempty" json:"type,omitempty" binding:"omitempty,oneof=1 2 3"`
	Content string   `bson:"content" json:"content,omitempty" binding:"omitempty"`
}
