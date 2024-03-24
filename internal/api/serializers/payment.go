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

type PaymentQRTopUpBodyValidate struct {
	Amount float64            `json:"amount" validate:"required"`
	Method primitive.ObjectID `json:"method" validate:"required"`
}

func (v *PaymentQRTopUpBodyValidate) Validate() error {
	return validator.Validate(v)
}

type PaymentQRTopUpResponse struct {
	Amount        float64 `json:"amount"`
	AccountName   string  `json:"account_name"`
	AccountNumber string  `json:"account_number"`
	TransactionId string  `json:"trx_id"`
	MethodCode    string  `json:"code"`
}

type PaymentCreateBodyValidate struct {
	Amount float64            `json:"amount"`
	UserId primitive.ObjectID `json:"user_id"`
	Method primitive.ObjectID `json:"method"`
}

func (v *PaymentCreateBodyValidate) Validate() error {
	return validator.Validate(v)
}

type PaymentListBodyValidate struct {
	Page     int64               `json:"page" validate:"omitempty"`
	Limit    int64               `json:"limit" validate:"omitempty"`
	SortType int                 `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string              `json:"sort_by" validate:"omitempty"`
	Search   string              `json:"search" validate:"omitempty"`
	Filters  []PaymentListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *PaymentListBodyValidate) Validate() error {
	return validator.Validate(v)
}

func (v *PaymentListBodyValidate) Sort() map[string]int {
	sortBy := "created_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type PaymentListFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq $regex"`
	Field  string      `json:"field" validate:"required,oneof=trx_id user_id"`
}

func (v *PaymentListBodyValidate) GetFilter() bson.M {
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

type PaymentListResponse struct {
	CreatedAt     time.Time          `json:"created_at"`
	Amount        float64            `json:"amount"`
	Status        int                `json:"status"`
	Feedback      string             `json:"feedback"`
	TransactionId string             `json:"trx_id"`
	Username      string             `json:"username"`
	Method        string             `json:"method"`
	Id            primitive.ObjectID `json:"id"`
}

type PaymentUpdateBodyValidate struct {
	Id       primitive.ObjectID `json:"id" validate:"required"`
	Status   int                `json:"status" validate:"required"`
	Feedback string             `json:"feedback" validate:"omitempty"`
}

func (v *PaymentUpdateBodyValidate) Validate() error {
	return validator.Validate(v)
}

type PaymentGetResponse struct {
	CreatedAt     time.Time          `json:"created_at"`
	Amount        float64            `json:"amount"`
	Status        int                `json:"status"`
	Feedback      string             `json:"feedback"`
	TransactionId string             `json:"trx_id"`
	Username      string             `json:"username"`
	Method        string             `json:"method"`
	Id            primitive.ObjectID `json:"id"`
}

type PaymentListByUserBodyValidate struct {
	Page     int64                     `json:"page" validate:"omitempty"`
	Limit    int64                     `json:"limit" validate:"omitempty"`
	SortType int                       `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string                    `json:"sort_by" validate:"omitempty"`
	Search   string                    `json:"search" validate:"omitempty"`
	Filters  []PaymentListByUserFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *PaymentListByUserBodyValidate) Validate() error {
	return validator.Validate(v)
}

func (v *PaymentListByUserBodyValidate) Sort() map[string]int {
	sortBy := "created_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type PaymentListByUserFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq $regex"`
	Field  string      `json:"field" validate:"required,oneof=trx_id"`
}

func (v *PaymentListByUserBodyValidate) GetFilter() bson.M {
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

type PaymentListByUserResponse struct {
	CreatedAt     time.Time          `json:"created_at"`
	Amount        float64            `json:"amount"`
	Status        int                `json:"status"`
	TransactionId string             `json:"trx_id"`
	Method        string             `json:"method"`
	Id            primitive.ObjectID `json:"id"`
}
