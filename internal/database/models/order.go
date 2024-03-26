package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	UpdatedAt             time.Time          `bson:"updated_at"`
	CreatedAt             time.Time          `bson:"created_at"`
	Quantity              int64              `bson:"quantity"`
	ProviderOrderId       int64              `bson:"provider_order_id"`
	ProviderOrderResponse string             `bson:"provider_order_response"`
	Status                int                `bson:"status"`
	StartCounter          int64              `bson:"start_counter"`
	Remains               int64              `bson:"remains"`
	Price                 float64            `bson:"price"`
	Note                  string             `bson:"note"`
	Link                  string             `bson:"link"`
	UserId                primitive.ObjectID `bson:"user_id"`
	ServiceId             primitive.ObjectID `bson:"service_id"`
	CategoryId            primitive.ObjectID `bson:"category_id"`
	ProviderId            primitive.ObjectID `bson:"provider_id"`
	Id                    primitive.ObjectID `bson:"_id,omitempty"`
}

const OrderCollectionName = "orders"
