package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	UpdatedAt     time.Time          `bson:"updated_at"`
	CreatedAt     time.Time          `bson:"created_at"`
	TransactionId string             `bson:"trx_id"`
	Amount        float64            `bson:"amount"`
	Type          int                `bson:"type"`
	UserId        primitive.ObjectID `bson:"user_id"`
	Id            primitive.ObjectID `bson:"_id,omitempty"`
}

const TransactionCollectionName = "transactions"
