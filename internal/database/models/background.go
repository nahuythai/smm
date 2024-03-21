package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Background struct {
	UpdatedAt time.Time          `bson:"updated_at"`
	CreatedAt time.Time          `bson:"created_at"`
	Type      int                `bson:"type"`
	UserId    primitive.ObjectID `bson:"user_id"`
	Id        primitive.ObjectID `bson:"_id,omitempty"`
}

const BackgroundCollectionName = "backgrounds"
