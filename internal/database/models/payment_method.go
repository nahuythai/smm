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
	MinAmount        int64              `bson:"min_amount"`
	MaxAmount        int64              `bson:"max_amount"`
	PercentageCharge float64            `bson:"percentage_charge"`
	FixedCharge      float64            `bson:"fixed_charge"`
	ConventionRate   float64            `bson:"convention_rate"`
	Description      string             `bson:"description"`
	Status           int                `bson:"status"`
	Auto             bool               `bson:"auto"`
	Extras map[string]interface{} `bson:"extras"`
	Id               primitive.ObjectID `bson:"_id,omitempty"`
}

const PaymentMethodCollectionName = "payment_methods"
