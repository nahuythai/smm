package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	UpdatedAt   time.Time          `bson:"updated_at"`
	CreatedAt   time.Time          `bson:"created_at"`
	Status      int                `bson:"status"`
	Title       string             `bson:"title"`
	Image       string             `bson:"image"`
	Description string             `bson:"description"`
	Id          primitive.ObjectID `bson:"_id,omitempty"`
}

const CategoryCollectionName = "categories"
