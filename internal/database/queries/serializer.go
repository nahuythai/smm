package queries

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	DisplayName       *string   `bson:"display_name,omitempty"`
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

type ServiceUpdateByIdDoc struct {
	UpdatedAt         time.Time          `bson:"updated_at"`
	Title             string             `bson:"title"`
	Status            int                `bson:"status"`
	MinAmount         int64              `bson:"min_amount"`
	MaxAmount         int64              `bson:"max_amount"`
	Rate              float64            `bson:"rate"`
	Description       string             `bson:"description"`
	ProviderServiceId string             `bson:"provider_service_id"`
	CategoryId        primitive.ObjectID `bson:"category_id"`
	ProviderId        primitive.ObjectID `bson:"provider_id"`
}

type ProviderUpdateByIdDoc struct {
	UpdatedAt   time.Time `bson:"updated_at"`
	ApiName     string    `bson:"api_name"`
	ApiKey      string    `bson:"api_key"`
	Description string    `bson:"description"`
	Url         string    `bson:"url"`
	Status      int       `bson:"status"`
	Rate        float64   `bson:"rate"`
}
