package serializers

import (
	"fmt"
	"regexp"
	"smm/internal/database/queries"
	"smm/pkg/validator"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserCreateBodyValidate struct {
	FirstName         string  `json:"first_name" validate:"omitempty"`
	LastName          string  `json:"last_name" validate:"omitempty"`
	Username          string  `json:"username" validate:"required,lowercase,alphanum,min=3,max=20"`
	PhoneNumber       string  `json:"phone_number" validate:"omitempty"`
	Language          int     `json:"language" validate:"omitempty"`
	Status            int     `json:"status" validate:"omitempty,oneof=0 1"`
	Balance           float64 `json:"balance" validate:"omitempty"`
	Email             string  `json:"email" validate:"required,lowercase,email"`
	EmailVerification bool    `json:"email_verification" validate:"omitempty"`
	TwoFAEnable       bool    `json:"2fa_enable" validate:"omitempty"`
	Address           string  `json:"address" validate:"omitempty"`
	Password          string  `json:"password" validate:"required,min=3,max=20"`
	Avatar            string  `json:"avatar" validate:"omitempty"`
}

func (v *UserCreateBodyValidate) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}
	return nil
}

type UserListBodyValidate struct {
	Page     int64            `json:"page" validate:"omitempty"`
	Limit    int64            `json:"limit" validate:"omitempty"`
	SortType int              `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string           `json:"sort_by" validate:"omitempty"`
	Search   string           `json:"search" validate:"omitempty"`
	Filters  []UserListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *UserListBodyValidate) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}
	return nil
}

func (v *UserListBodyValidate) Sort() map[string]int {
	sortBy := "updated_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type UserListFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq $regex"`
	Field  string      `json:"field" validate:"required,oneof=username email phone_number"`
}

func (v *UserListBodyValidate) GetFilter() bson.M {
	filters := make([]queries.Filter, 0, len(v.Filters))
	for _, filter := range v.Filters {
		filters = append(filters, queries.Filter(filter))
	}
	filterOption := queries.NewFilterOption()
	filterOption.AddFilter(filters...)
	query := filterOption.BuildAndQuery()
	if v.Search != "" {
		query["$or"] = []bson.M{{"username": bson.M{queries.QueryFilterMethodRegex: primitive.Regex{Pattern: regexp.QuoteMeta(fmt.Sprintf("%v", v.Search)), Options: "i"}}}}
	}
	return query
}

type UserListResponse struct {
	UpdatedAt         time.Time          `json:"updated_at"`
	CreatedAt         time.Time          `json:"created_at"`
	LastActive        *time.Time         `json:"last_active"`
	FirstName         string             `json:"first_name"`
	LastName          string             `json:"last_name"`
	Username          string             `json:"username"`
	PhoneNumber       string             `json:"phone_number"`
	Language          int                `json:"language"`
	Status            int                `json:"status"`
	Balance           float64            `json:"balance"`
	Email             string             `json:"email"`
	EmailVerification bool               `json:"email_verification"`
	TwoFAEnable       bool               `json:"2fa_enable"`
	Address           string             `json:"address"`
	Avatar            string             `json:"avatar"`
	ApiToken          string             `json:"api_token"`
	Id                primitive.ObjectID `json:"id"`
}
