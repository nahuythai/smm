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

type PaymentTopUpBodyValidate struct {
	Amount        float64            `json:"amount"`
	PaymentMethod primitive.ObjectID `json:"payment_method"`
}

func (v *PaymentTopUpBodyValidate) Validate() error {
	return validator.Validate(v)
}

type PaymentCreateBodyValidate struct {
	Amount        float64            `json:"amount"`
	UserId        primitive.ObjectID `json:"user_id"`
	PaymentMethod primitive.ObjectID `json:"payment_method"`
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
	Field  string      `json:"field" validate:"required,oneof=trx_id"`
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
	CreatedAt     time.Time              `json:"created_at"`
	Amount        float64                `json:"amount"`
	Status        int                    `json:"status"`
	Feedback      string                 `json:"feedback"`
	TransactionId string                 `json:"trx_id"`
	Username      string                 `json:"username"`
	PaymentMethod map[string]interface{} `json:"payment_method"`
	Id            primitive.ObjectID     `json:"_id,omitempty"`
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
	CreatedAt     time.Time              `json:"created_at"`
	Amount        float64                `json:"amount"`
	Status        int                    `json:"status"`
	Feedback      string                 `json:"feedback"`
	TransactionId string                 `json:"trx_id"`
	Username      string                 `json:"username"`
	PaymentMethod map[string]interface{} `json:"payment_method"`
	Id            primitive.ObjectID     `json:"id"`
}
