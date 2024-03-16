package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Provider struct {
	UpdatedAt   time.Time          `bson:"updated_at"`
	CreatedAt   time.Time          `bson:"created_at"`
	ApiName     string             `bson:"api_name"`
	ApiKey      string             `bson:"api_key"`
	Url         string             `bson:"url"`
	Description string             `bson:"description"`
	Status      int                `bson:"status"`
	Rate        float64            `bson:"rate"`
	Balance     float64            `bson:"balance"`
	Id          primitive.ObjectID `bson:"_id,omitempty"`
}

const ProviderCollectionName = "providers"
