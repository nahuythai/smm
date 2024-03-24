package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentMethod struct {
	UpdatedAt        time.Time          `bson:"updated_at"`
	CreatedAt        time.Time          `bson:"created_at"`
	Name             string             `bson:"name"`
	Code             string             `bson:"code"`
	Image            string             `bson:"image"`
	MinAmount        float64            `bson:"min_amount"`
	MaxAmount        float64            `bson:"max_amount"`
	PercentageCharge float64            `bson:"percentage_charge"`
	FixedCharge      float64            `bson:"fixed_charge"`
	ConventionRate   float64            `bson:"convention_rate"`
	Description      string             `bson:"description"`
	Status           int                `bson:"status"`
	Auto             bool               `bson:"auto"`
	AccountName      string             `bson:"account_name"`
	AccountNumber    string             `bson:"account_number"`
	Id               primitive.ObjectID `bson:"_id,omitempty"`
}

const PaymentMethodCollectionName = "payment_methods"
