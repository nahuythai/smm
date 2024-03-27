package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	UpdatedAt         time.Time          `bson:"updated_at"`
	CreatedAt         time.Time          `bson:"created_at"`
	Title             string             `bson:"title"`
	Status            int                `bson:"status"`
	MinAmount         int64              `bson:"min_amount"`
	MaxAmount         int64              `bson:"max_amount"`
	Rate              float64            `bson:"rate"`
	Description       string             `bson:"description"`
	ProviderServiceId string             `bson:"provider_service_id"`
	CategoryId        primitive.ObjectID `bson:"category_id"`
	ProviderId        primitive.ObjectID `bson:"provider_id"`
	Seq               int                `bson:"seq"`
	Id                primitive.ObjectID `bson:"_id,omitempty"`
}

const ServiceCollectionName = "services"
