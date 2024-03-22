package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	UpdatedAt     time.Time          `bson:"updated_at"`
	CreatedAt     time.Time          `bson:"created_at"`
	Amount        float64            `bson:"amount"`
	Status        int                `bson:"status"`
	Feedback      string             `bson:"feedback"`
	UserId        primitive.ObjectID `bson:"user_id"`
	PaymentMethod primitive.ObjectID `bson:"payment_method"`
	Id            primitive.ObjectID `bson:"_id,omitempty"`
}

const PaymentCollectionName = "payment_methods"
