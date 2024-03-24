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

type OrderUpdateByIdDoc struct {
	UpdatedAt    time.Time `bson:"updated_at"`
	Status       int       `bson:"status"`
	StartCounter int64     `bson:"start_counter"`
	Remains      int64     `bson:"remains"`
	Note         string    `bson:"note"`
	Link         string    `bson:"link"`
}

type PaymentMethodUpdateByIdDoc struct {
	UpdatedAt        time.Time `bson:"updated_at"`
	Name             string    `bson:"name"`
	Code             string    `bson:"code"`
	Image            string    `bson:"image"`
	MinAmount        float64   `bson:"min_amount"`
	MaxAmount        float64   `bson:"max_amount"`
	Status           int       `bson:"status"`
	PercentageCharge float64   `bson:"percentage_charge"`
	Description      string    `bson:"description"`
	Auto             bool      `bson:"auto"`
	AccountName      string    `bson:"account_name"`
	AccountNumber    string    `bson:"account_number"`
}

type PaymentUpdateByIdDoc struct {
	UpdatedAt time.Time `bson:"updated_at"`
	Status    int       `bson:"status"`
	Feedback  string    `bson:"feedback"`
}
