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
	DisplayName       string  `json:"display_name" validate:"omitempty"`
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
	Username    string             `json:"username"`
	DisplayName string             `json:"display_name"`
	PhoneNumber string             `json:"phone_number"`
	Status      int                `json:"status"`
	Balance     float64            `json:"balance"`
	Email       string             `json:"email"`
	Avatar      string             `json:"avatar"`
	Id          primitive.ObjectID `json:"id"`
}

type UserGenerateAPIKey struct {
	Id primitive.ObjectID `json:"id"`
}

func (v *UserGenerateAPIKey) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}
	return nil
}

type UserUpdateBodyValidate struct {
	Id                primitive.ObjectID `json:"id" validate:"required"`
	DisplayName       *string            `json:"display_name" validate:"omitempty"`
	Username          string             `json:"username" validate:"required,lowercase,alphanum,min=3,max=20"`
	PhoneNumber       *string            `json:"phone_number" validate:"omitempty"`
	Language          *int               `json:"language" validate:"omitempty"`
	Status            *int               `json:"status" validate:"omitempty,oneof=0 1"`
	Email             string             `json:"email" validate:"required,lowercase,email"`
	EmailVerification *bool              `json:"email_verification" validate:"omitempty"`
	TwoFAEnable       *bool              `json:"2fa_enable" validate:"omitempty"`
	Address           *string            `json:"address" validate:"omitempty"`
	Avatar            *string            `json:"avatar" validate:"omitempty"`
}

func (v *UserUpdateBodyValidate) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}
	return nil
}

type UserGetResponse struct {
	UpdatedAt         time.Time          `json:"updated_at"`
	CreatedAt         time.Time          `json:"created_at"`
	LastActive        *time.Time         `json:"last_active"`
	DisplayName       string             `json:"display_name"`
	Username          string             `json:"username"`
	PhoneNumber       string             `json:"phone_number"`
	Language          int                `json:"language"`
	Status            int                `json:"status"`
	Role              int                `json:"role"`
	Balance           float64            `json:"balance"`
	Email             string             `json:"email"`
	EmailVerification bool               `json:"email_verification"`
	TwoFAEnable       bool               `json:"2fa_enable"`
	Address           string             `json:"address"`
	Avatar            string             `json:"avatar"`
	ApiKey            string             `json:"api_key"`
	Id                primitive.ObjectID `json:"id"`
}

type UserUpdateBalanceBodyValidate struct {
	Id     primitive.ObjectID `json:"id" validate:"required"`
	Amount float64            `json:"amount" validate:"required"`
}

func (v *UserUpdateBalanceBodyValidate) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}
	return nil
}

type UserUpdatePasswordBodyValidate struct {
	Id       primitive.ObjectID `json:"id" validate:"required"`
	Password string             `json:"password" validate:"required,min=3,max=20"`
}

func (v *UserUpdatePasswordBodyValidate) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}
	return nil
}

type UserLoginBodyValidate struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=3,max=20"`
}

func (v *UserLoginBodyValidate) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}
	return nil
}

type UserLoginVerifyBodyValidate struct {
	Code string `json:"code" validate:"required"`
}

func (v *UserLoginVerifyBodyValidate) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}
	return nil
}

type UserRegisterBodyValidate struct {
	DisplayName string `json:"display_name" validate:"omitempty"`
	Username    string `json:"username" validate:"required,lowercase,alphanum,min=3,max=20"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email" validate:"required,lowercase,email"`
	Password    string `json:"password" validate:"required,min=3,max=20"`
}

func (v *UserRegisterBodyValidate) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}
	return nil
}
