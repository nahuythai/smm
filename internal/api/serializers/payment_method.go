package serializers

import (
	"fmt"
	"regexp"
	"smm/internal/database/queries"
	"smm/pkg/validator"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentMethodCreateBodyValidate struct {
	Name             string  `json:"name" validate:"required"`
	Code             string  `json:"code" validate:"required"`
	Image            string  `json:"image" validate:"omitempty"`
	MinAmount        float64 `json:"min_amount" validate:"required,min=1"`
	MaxAmount        float64 `json:"max_amount" validate:"required,min=1"`
	PercentageCharge float64 `json:"percentage_charge" validate:"omitempty,min=0"`
	FixedCharge      float64 `json:"fixed_charge" validate:"omitempty,min=0"`
	ConventionRate   float64 `json:"convention_rate" validate:"omitempty,min=0"`
	Description      string  `json:"description" validate:"omitempty"`
	Auto             *bool   `json:"auto" validate:"required"`
	AccountName      string  `json:"account_name" validate:"required"`
	AccountNumber    string  `json:"account_number" validate:"required"`
}

func (v *PaymentMethodCreateBodyValidate) Validate() error {
	return validator.Validate(v)
}

type PaymentMethodListBodyValidate struct {
	Page     int64                     `json:"page" validate:"omitempty"`
	Limit    int64                     `json:"limit" validate:"omitempty"`
	SortType int                       `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string                    `json:"sort_by" validate:"omitempty"`
	Search   string                    `json:"search" validate:"omitempty"`
	Filters  []PaymentMethodListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *PaymentMethodListBodyValidate) Validate() error {
	return validator.Validate(v)
}

func (v *PaymentMethodListBodyValidate) Sort() map[string]int {
	sortBy := "created_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type PaymentMethodListFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq $regex"`
	Field  string      `json:"field" validate:"required,oneof=name code"`
}

func (v *PaymentMethodListBodyValidate) GetFilter() bson.M {
	filters := make([]queries.Filter, 0, len(v.Filters))
	for _, filter := range v.Filters {
		filters = append(filters, queries.Filter(filter))
	}
	filterOption := queries.NewFilterOption()
	filterOption.AddFilter(filters...)
	query := filterOption.BuildAndQuery()
	if v.Search != "" {
		query["$or"] = []bson.M{{"title": bson.M{queries.QueryFilterMethodRegex: primitive.Regex{Pattern: regexp.QuoteMeta(fmt.Sprintf("%v", v.Search)), Options: "i"}}}}
	}
	return query
}

type PaymentMethodListResponse struct {
	Name   string             `json:"name"`
	Code   string             `json:"code"`
	Image  string             `json:"image"`
	Status int                `json:"status"`
	Id     primitive.ObjectID `json:"id"`
}

type PaymentMethodUpdateBodyValidate struct {
	Id               primitive.ObjectID `json:"id" validate:"required"`
	Name             string             `json:"name" validate:"required"`
	Code             string             `json:"code" validate:"required"`
	Image            string             `json:"image" validate:"omitempty"`
	MinAmount        float64            `json:"min_amount" validate:"required,min=1"`
	MaxAmount        float64            `json:"max_amount" validate:"required,min=1"`
	Status           int                `json:"status" validate:"omitempty,oneof=0 1"`
	PercentageCharge float64            `json:"percentage_charge" validate:"omitempty,min=0"`
	FixedCharge      float64            `json:"fixed_charge" validate:"omitempty,min=0"`
	ConventionRate   float64            `json:"convention_rate" validate:"omitempty,min=0"`
	Description      string             `json:"description" validate:"omitempty"`
	Auto             *bool              `json:"auto" validate:"required"`
	AccountName      string             `json:"account_name" validate:"required"`
	AccountNumber    string             `json:"account_number" validate:"required"`
}

func (v *PaymentMethodUpdateBodyValidate) Validate() error {
	return validator.Validate(v)
}

type PaymentMethodGetResponse struct {
	Name             string             `json:"name"`
	Code             string             `json:"code"`
	Image            string             `json:"image"`
	MinAmount        float64            `json:"min_amount"`
	MaxAmount        float64            `json:"max_amount"`
	PercentageCharge float64            `json:"percentage_charge"`
	FixedCharge      float64            `json:"fixed_charge"`
	ConventionRate   float64            `json:"convention_rate"`
	Description      string             `json:"description"`
	Status           int                `json:"status"`
	Auto             bool               `json:"auto"`
	AccountName      string             `json:"account_name"`
	AccountNumber    string             `json:"account_number"`
	Id               primitive.ObjectID `json:"id"`
}
