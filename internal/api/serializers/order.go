package serializers

import (
	"smm/internal/database/queries"
	"smm/pkg/validator"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderCreateBodyValidate struct {
	Quantity  int64              `json:"quantity" validate:"required"`
	Link      string             `json:"link" validate:"required"`
	ServiceId primitive.ObjectID `json:"service_id" validate:"required"`
}

func (v *OrderCreateBodyValidate) Validate() error {
	return validator.Validate(v)
}

type OrderListBodyValidate struct {
	Page     int64             `json:"page" validate:"omitempty"`
	Limit    int64             `json:"limit" validate:"omitempty"`
	SortType int               `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string            `json:"sort_by" validate:"omitempty"`
	Search   string            `json:"search" validate:"omitempty"`
	Filters  []OrderListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *OrderListBodyValidate) Validate() error {
	return validator.Validate(v)
}

func (v *OrderListBodyValidate) Sort() map[string]int {
	sortBy := "created_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type OrderListFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq"`
	Field  string      `json:"field" validate:"required,oneof=status user_id"`
}

func (v *OrderListBodyValidate) GetFilter() bson.M {
	filters := make([]queries.Filter, 0, len(v.Filters))
	for _, filter := range v.Filters {
		filters = append(filters, queries.Filter(filter))
	}
	filterOption := queries.NewFilterOption()
	filterOption.AddFilter(filters...)
	query := filterOption.BuildAndQuery()
	// if v.Search != "" {
	// 	// query["$or"] = []bson.M{{"title": bson.M{queries.QueryFilterMethodRegex: primitive.Regex{Pattern: regexp.QuoteMeta(fmt.Sprintf("%v", v.Search)), Options: "i"}}}}
	// }
	return query
}

type OrderListResponse struct {
	UpdatedAt    time.Time          `json:"updated_at"`
	Quantity     int64              `json:"quantity"`
	Status       int                `json:"status"`
	StartCounter int64              `json:"start_counter"`
	Remains      int64              `json:"remains"`
	Price        float64            `json:"price"`
	Link         string             `json:"link"`
	Service      string             `json:"service"`
	Username     string             `json:"username"`
	Id           primitive.ObjectID `json:"id"`
}

type OrderListByUserBodyValidate struct {
	Page     int64             `json:"page" validate:"omitempty"`
	Limit    int64             `json:"limit" validate:"omitempty"`
	SortType int               `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string            `json:"sort_by" validate:"omitempty"`
	Search   string            `json:"search" validate:"omitempty"`
	Filters  []OrderListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *OrderListByUserBodyValidate) Validate() error {
	return validator.Validate(v)
}

func (v *OrderListByUserBodyValidate) Sort() map[string]int {
	sortBy := "created_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type OrderListByUserFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq"`
	Field  string      `json:"field" validate:"required,oneof=status"`
}

func (v *OrderListByUserBodyValidate) GetFilter() bson.M {
	filters := make([]queries.Filter, 0, len(v.Filters))
	for _, filter := range v.Filters {
		filters = append(filters, queries.Filter(filter))
	}
	filterOption := queries.NewFilterOption()
	filterOption.AddFilter(filters...)
	query := filterOption.BuildAndQuery()
	// if v.Search != "" {
	// 	// query["$or"] = []bson.M{{"title": bson.M{queries.QueryFilterMethodRegex: primitive.Regex{Pattern: regexp.QuoteMeta(fmt.Sprintf("%v", v.Search)), Options: "i"}}}}
	// }
	return query
}

type OrderListByUserResponse struct {
	CreatedAt    time.Time          `json:"created_at"`
	Quantity     int64              `json:"quantity"`
	Status       int                `json:"status"`
	StartCounter int64              `json:"start_counter"`
	Remains      int64              `json:"remains"`
	Price        float64            `json:"price"`
	Link         string             `json:"link"`
	Service      string             `json:"service"`
	Id           primitive.ObjectID `json:"id"`
}

type OrderUpdateBodyValidate struct {
	Status       int                `json:"status" validate:"omitempty"`
	StartCounter int64              `json:"start_counter" validate:"omitempty"`
	Remains      int64              `json:"remains" validate:"omitempty"`
	Note         string             `json:"note" validate:"omitempty"`
	Link         string             `json:"link" validate:"required"`
	Id           primitive.ObjectID `json:"id" validate:"required"`
}

func (v *OrderUpdateBodyValidate) Validate() error {
	return validator.Validate(v)
}

type OrderGetResponse struct {
	Quantity              int64              `json:"quantity"`
	ProviderOrderId       int64              `json:"provider_order_id"`
	ProviderOrderResponse string             `json:"provider_order_response"`
	Status                int                `json:"status"`
	StartCounter          int64              `json:"start_counter"`
	Remains               int64              `json:"remains"`
	Note                  string             `json:"note"`
	Link                  string             `json:"link"`
	Service               string             `json:"service"`
	Category              string             `json:"category"`
	Provider              string             `json:"provider"`
	Id                    primitive.ObjectID `json:"id"`
}
