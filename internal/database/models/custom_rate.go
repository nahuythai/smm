package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomRate struct {
	UpdatedAt  time.Time          `bson:"updated_at"`
	CreatedAt  time.Time          `bson:"created_at"`
	Price      float64            `bson:"price"`
	UserId     primitive.ObjectID `bson:"title"`
	ServiceId  primitive.ObjectID `bson:"service_id"`
	CategoryId primitive.ObjectID `bson:"category_id"`
	Id         primitive.ObjectID `bson:"_id,omitempty"`
}

const CustomRateCollectionName = "custom_rates"
