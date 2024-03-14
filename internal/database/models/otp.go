package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Otp struct {
	UpdatedAt time.Time          `bson:"updated_at"`
	CreatedAt time.Time          `bson:"created_at"`
	ExpiredAt time.Time          `bson:"expired_at"`
	Code      string             `bson:"code"`
	SecretKey string             `bson:"secret_key"`
	OtpUrl    string             `bson:"otp_url"`
	UserId    primitive.ObjectID `bson:"user_id"`
	Id        primitive.ObjectID `bson:"_id,omitempty"`
}

const OtpCollectionName = "otps"
