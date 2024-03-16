package serializers

import (
	"fmt"
	"regexp"
	"smm/internal/database/queries"
	"smm/pkg/validator"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProviderCreateBodyValidate struct {
	ApiName     string  `json:"api_name" validate:"required"`
	ApiKey      string  `json:"api_key" validate:"required"`
	Url         string  `json:"url" validate:"required"`
	Description string  `json:"description" validate:"omitempty"`
	Rate        float64 `json:"rate" validate:"required"`
}

func (v *ProviderCreateBodyValidate) Validate() error {
	return validator.Validate(v)
}

type ProviderListBodyValidate struct {
	Page     int64                `json:"page" validate:"omitempty"`
	Limit    int64                `json:"limit" validate:"omitempty"`
	SortType int                  `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string               `json:"sort_by" validate:"omitempty"`
	Search   string               `json:"search" validate:"omitempty"`
	Filters  []ProviderListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *ProviderListBodyValidate) Validate() error {
	return validator.Validate(v)
}

func (v *ProviderListBodyValidate) Sort() map[string]int {
	sortBy := "updated_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type ProviderListFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq $regex"`
	Field  string      `json:"field" validate:"required,oneof=username email phone_number"`
}

func (v *ProviderListBodyValidate) GetFilter() bson.M {
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

type ProviderListResponse struct {
	ApiName     string             `json:"api_name"`
	Status      int                `json:"status"`
	Balance     float64            `json:"balance"`
	Description string             `json:"description"`
	Id          primitive.ObjectID `json:"id"`
}

type ProviderUpdateBodyValidate struct {
	Id          primitive.ObjectID `json:"id" validate:"required"`
	ApiName     string             `json:"api_name" validate:"required"`
	Url         string             `json:"url" validate:"required"`
	ApiKey      string             `json:"api_key" validate:"required"`
	Description string             `json:"description" validate:"omitempty"`
	Status      *int               `json:"status" validate:"required"`
	Rate        float64            `json:"rate" validate:"required"`
}

func (v *ProviderUpdateBodyValidate) Validate() error {
	return validator.Validate(v)
}

type ProviderGetResponse struct {
	Status      int                `json:"status"`
	ApiName     string             `json:"api_name"`
	ApiKey      string             `json:"api_key"`
	Description string             `json:"description"`
	Url         string             `json:"url"`
	Balance     float64            `json:"balance"`
	Rate        float64            `json:"rate"`
	Id          primitive.ObjectID `json:"id"`
}
