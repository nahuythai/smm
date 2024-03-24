package serializers

import (
	"fmt"
	"regexp"
	"smm/internal/database/queries"
	"smm/pkg/validator"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceCreateBodyValidate struct {
	Title             string             `json:"title" validate:"required"`
	MinAmount         int64              `json:"min_amount" validate:"required,min=1"`
	MaxAmount         int64              `json:"max_amount" validate:"required,min=1"`
	Rate              float64            `json:"rate" validate:"required,min=0"`
	Description       string             `json:"description" validate:"omitempty"`
	ProviderServiceId string             `json:"provider_service_id" validate:"required"`
	CategoryId        primitive.ObjectID `json:"category_id" validate:"required"`
	ProviderId        primitive.ObjectID `json:"provider_id" validate:"required"`
}

func (v *ServiceCreateBodyValidate) Validate() error {
	return validator.Validate(v)
}

type ServiceListBodyValidate struct {
	Page     int64               `json:"page" validate:"omitempty"`
	Limit    int64               `json:"limit" validate:"omitempty"`
	SortType int                 `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string              `json:"sort_by" validate:"omitempty"`
	Search   string              `json:"search" validate:"omitempty"`
	Filters  []ServiceListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *ServiceListBodyValidate) Validate() error {
	return validator.Validate(v)
}

func (v *ServiceListBodyValidate) Sort() map[string]int {
	sortBy := "updated_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type ServiceListFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq $regex"`
	Field  string      `json:"field" validate:"required,oneof=title category_id status provider_id"`
}

func (v *ServiceListBodyValidate) GetFilter() bson.M {
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

type ServiceListResponse struct {
	Title       string             `json:"title"`
	Status      int                `json:"status"`
	Provider    string             `json:"provider"`
	MinAmount   int64              `json:"min_amount"`
	MaxAmount   int64              `json:"max_amount"`
	Rate        float64            `json:"rate"`
	Description string             `json:"description"`
	Id          primitive.ObjectID `json:"id"`
}

type ServiceUpdateBodyValidate struct {
	Id                primitive.ObjectID `json:"id" validate:"required"`
	Title             string             `json:"title" validate:"required"`
	Status            int                `json:"status" validate:"oneof=0 1"`
	MinAmount         int64              `json:"min_amount" validate:"required,min=1"`
	MaxAmount         int64              `json:"max_amount" validate:"required,max=1"`
	Rate              float64            `json:"rate" validate:"required,min=0"`
	Description       string             `json:"description" validate:"omitempty"`
	ProviderServiceId string             `json:"provider_service_id" validate:"required"`
	ProviderId        primitive.ObjectID `json:"provider_id" validate:"required"`
	CategoryId        primitive.ObjectID `json:"category_id" validate:"required"`
}

func (v *ServiceUpdateBodyValidate) Validate() error {
	return validator.Validate(v)
}

type ServiceGetResponse struct {
	Id                primitive.ObjectID `json:"id"`
	Title             string             `json:"title"`
	Status            int                `json:"status"`
	MinAmount         int64              `json:"min_amount"`
	MaxAmount         int64              `json:"max_amount"`
	Rate              float64            `json:"rate"`
	Description       string             `json:"description"`
	ProviderServiceId string             `json:"provider_service_id"`
	Provider          string             `json:"provider"`
	Category          string             `json:"category"`
}

type ServiceUserListBodyValidate struct {
	Page     int64                   `json:"page" validate:"omitempty"`
	Limit    int64                   `json:"limit" validate:"omitempty"`
	SortType int                     `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string                  `json:"sort_by" validate:"omitempty"`
	Search   string                  `json:"search" validate:"omitempty"`
	Filters  []ServiceUserListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *ServiceUserListBodyValidate) Validate() error {
	return validator.Validate(v)
}

func (v *ServiceUserListBodyValidate) Sort() map[string]int {
	sortBy := "updated_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type ServiceUserListFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq $regex"`
	Field  string      `json:"field" validate:"required,oneof=title category_id provider_id"`
}

func (v *ServiceUserListBodyValidate) GetFilter() bson.M {
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

type ServiceUserListResponse struct {
	Title       string             `json:"title"`
	Status      int                `json:"status"`
	Provider    string             `json:"provider"`
	MinAmount   int64              `json:"min_amount"`
	MaxAmount   int64              `json:"max_amount"`
	Rate        float64            `json:"rate"`
	Description string             `json:"description"`
	Id          primitive.ObjectID `json:"id"`
}
