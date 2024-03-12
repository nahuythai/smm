package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UpdatedAt         time.Time          `bson:"updated_at"`
	CreatedAt         time.Time          `bson:"created_at"`
	LastActive        *time.Time         `bson:"last_active"`
	FirstName         string             `bson:"first_name"`
	LastName          string             `bson:"last_name"`
	Username          string             `bson:"username"`
	PhoneNumber       string             `bson:"phone_number"`
	Language          int                `bson:"language"`
	Status            int                `bson:"status"`
	Balance           float64            `bson:"balance"`
	Email             string             `bson:"email"`
	EmailVerification bool               `bson:"email_verification"`
	TwoFAEnable       bool               `bson:"2fa_enable"`
	Address           string             `bson:"address"`
	Password          string             `bson:"password"`
	Avatar            string             `bson:"avatar"`
	ApiKey            string             `bson:"api_key"`
	Id                primitive.ObjectID `bson:"_id,omitempty"`
}

const UserCollectionName = "users"
