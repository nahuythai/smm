package queries

import (
	"time"
)

type CategoryUpdateByIdDoc struct {
	UpdatedAt   time.Time `bson:"updated_at"`
	Title       string    `bson:"title"`
	Image       *string   `bson:"image,omitempty"`
	Status      *int      `bson:"status,omitempty"`
	Description *string   `bson:"description,omitempty"`
}

type UserUpdateByIdDoc struct {
	UpdatedAt         time.Time `bson:"updated_at"`
	FirstName         *string   `bson:"first_name,omitempty"`
	LastName          *string   `bson:"last_name,omitempty"`
	Username          string    `bson:"username"`
	PhoneNumber       *string   `bson:"phone_number,omitempty"`
	Language          *int      `bson:"language,omitempty"`
	Status            *int      `bson:"status,omitempty"`
	Email             string    `bson:"email"`
	EmailVerification *bool     `bson:"email_verification,omitempty"`
	TwoFAEnable       *bool     `bson:"2fa_enable,omitempty"`
	Address           *string   `bson:"address,omitempty"`
	Avatar            *string   `bson:"avatar,omitempty"`
}
