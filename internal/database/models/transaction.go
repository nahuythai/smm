package models

import (
	"crypto/rand"
	"math/big"
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

func (t *Transaction) GenerateTransactionId() string {
	length := 10
	charSet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomString := make([]byte, length)
	maxIndex := big.NewInt(int64(len(charSet)))

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, maxIndex)
		if err != nil {
			return "----------"
		}
		randomString[i] = charSet[randomIndex.Int64()]
	}

	return string(randomString)
}
