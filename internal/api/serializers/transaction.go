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

type TransactionListBodyValidate struct {
	Page     int64                   `json:"page" validate:"omitempty"`
	Limit    int64                   `json:"limit" validate:"omitempty"`
	SortType int                     `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string                  `json:"sort_by" validate:"omitempty"`
	Search   string                  `json:"search" validate:"omitempty"`
	Filters  []TransactionListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *TransactionListBodyValidate) Validate() error {
	return validator.Validate(v)
}

func (v *TransactionListBodyValidate) Sort() map[string]int {
	sortBy := "created_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type TransactionListFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq $regex"`
	Field  string      `json:"field" validate:"required,oneof=trx_id"`
}

func (v *TransactionListBodyValidate) GetFilter() bson.M {
	filters := make([]queries.Filter, 0, len(v.Filters))
	for _, filter := range v.Filters {
		filters = append(filters, queries.Filter(filter))
	}
	filterOption := queries.NewFilterOption()
	filterOption.AddFilter(filters...)
	query := filterOption.BuildAndQuery()
	if v.Search != "" {
		query["$or"] = []bson.M{{"trx_id": bson.M{queries.QueryFilterMethodRegex: primitive.Regex{Pattern: regexp.QuoteMeta(fmt.Sprintf("%v", v.Search)), Options: "i"}}}}
	}
	return query
}

type TransactionListResponse struct {
	CreatedAt     time.Time `json:"created_at"`
	TransactionId string    `json:"trx_id"`
	Amount        float64   `json:"amount"`
	Type          int       `json:"type"`
	Username      string    `json:"username"`
}

type TransactionUserListBodyValidate struct {
	Page     int64                       `json:"page" validate:"omitempty"`
	Limit    int64                       `json:"limit" validate:"omitempty"`
	SortType int                         `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string                      `json:"sort_by" validate:"omitempty"`
	Search   string                      `json:"search" validate:"omitempty"`
	Filters  []TransactionUserListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *TransactionUserListBodyValidate) Validate() error {
	return validator.Validate(v)
}

func (v *TransactionUserListBodyValidate) Sort() map[string]int {
	sortBy := "created_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type TransactionUserListFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq $regex"`
	Field  string      `json:"field" validate:"required,oneof=trx_id"`
}

func (v *TransactionUserListBodyValidate) GetFilter() bson.M {
	filters := make([]queries.Filter, 0, len(v.Filters))
	for _, filter := range v.Filters {
		filters = append(filters, queries.Filter(filter))
	}
	filterOption := queries.NewFilterOption()
	filterOption.AddFilter(filters...)
	query := filterOption.BuildAndQuery()
	if v.Search != "" {
		query["$or"] = []bson.M{{"trx_id": bson.M{queries.QueryFilterMethodRegex: primitive.Regex{Pattern: regexp.QuoteMeta(fmt.Sprintf("%v", v.Search)), Options: "i"}}}}
	}
	return query
}

type TransactionUserListResponse struct {
	CreatedAt     time.Time `json:"created_at"`
	TransactionId string    `json:"trx_id"`
	Amount        float64   `json:"amount"`
	Type          int       `json:"type"`
}
