package serializers

import (
	"errors"
	"smm/internal/database/queries"
	"smm/pkg/constants"
	"smm/pkg/validator"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThirdPartyRouteValidate struct {
	Action string `json:"action" validate:"required"`
}

func (v *ThirdPartyRouteValidate) Validate() error {
	mapAction := map[string]struct{}{
		constants.ThirdPartyActionOrderCreate:         {},
		constants.ThirdPartyActionOrderStatus:         {},
		constants.ThirdPartyActionOrderMultipleStatus: {},
		constants.ThirdPartyActionServiceList:         {},
		constants.ThirdPartyActionUserBalance:         {},
	}
	if _, ok := mapAction[v.Action]; !ok {
		return errors.New("the selected action is invalid")
	}
	return validator.Validate(v)
}

type ThirdPartyOrderCreateBodyValidate struct {
	Quantity  int64  `json:"quantity" validate:"required"`
	Link      string `json:"link" validate:"required"`
	ServiceId string `json:"service_id" validate:"required"`
}

func (v *ThirdPartyOrderCreateBodyValidate) Validate() error {
	return validator.Validate(v)
}

type ThirdPartyOrderListBodyValidate struct {
	Page     int64                       `json:"page" validate:"omitempty"`
	Limit    int64                       `json:"limit" validate:"omitempty"`
	SortType int                         `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string                      `json:"sort_by" validate:"omitempty"`
	Search   string                      `json:"search" validate:"omitempty"`
	Filters  []ThirdPartyOrderListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *ThirdPartyOrderListBodyValidate) Validate() error {
	return validator.Validate(v)
}

func (v *ThirdPartyOrderListBodyValidate) Sort() map[string]int {
	sortBy := "created_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type ThirdPartyOrderListFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq"`
	Field  string      `json:"field" validate:"required,oneof=status user_id"`
}

func (v *ThirdPartyOrderListBodyValidate) GetFilter() bson.M {
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

type ThirdPartyOrderListResponse struct {
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

type ThirdPartyOrderListByUserBodyValidate struct {
	Page     int64                       `json:"page" validate:"omitempty"`
	Limit    int64                       `json:"limit" validate:"omitempty"`
	SortType int                         `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string                      `json:"sort_by" validate:"omitempty"`
	Search   string                      `json:"search" validate:"omitempty"`
	Filters  []ThirdPartyOrderListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *ThirdPartyOrderListByUserBodyValidate) Validate() error {
	return validator.Validate(v)
}

func (v *ThirdPartyOrderListByUserBodyValidate) Sort() map[string]int {
	sortBy := "created_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type ThirdPartyOrderListByUserFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq"`
	Field  string      `json:"field" validate:"required,oneof=status"`
}

func (v *ThirdPartyOrderListByUserBodyValidate) GetFilter() bson.M {
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

type ThirdPartyOrderListByUserResponse struct {
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

type ThirdPartyOrderUpdateBodyValidate struct {
	Status       int                `json:"status" validate:"omitempty"`
	StartCounter int64              `json:"start_counter" validate:"omitempty"`
	Remains      int64              `json:"remains" validate:"omitempty"`
	Note         string             `json:"note" validate:"omitempty"`
	Link         string             `json:"link" validate:"required"`
	Id           primitive.ObjectID `json:"id" validate:"required"`
}

func (v *ThirdPartyOrderUpdateBodyValidate) Validate() error {
	return validator.Validate(v)
}

type ThirdPartyOrderStatusResponse struct {
	Quantity     int64   `json:"quantity"`
	Status       string  `json:"status"`
	StartCounter int64   `json:"start_counter"`
	Remains      int64   `json:"remains"`
	Currency     string  `json:"currency"`
	Charge       *string `json:"charge"`
}

type ThirdPartyOrderStatusBodyValidate struct {
	OrderId string `json:"order" validate:"required"`
}

func (v *ThirdPartyOrderStatusBodyValidate) Validate() error {
	return validator.Validate(v)
}

type ThirdPartyMultipleOrderStatusBodyValidate struct {
	Orders string `json:"orders" validate:"required"`
}

func (v *ThirdPartyMultipleOrderStatusBodyValidate) Validate() error {
	return validator.Validate(v)
}

type ThirdPartyMultipleOrderStatusResponse struct {
	Quantity     int64              `json:"quantity"`
	Status       string             `json:"status"`
	StartCounter int64              `json:"start_counter"`
	Remains      int64              `json:"remains"`
	Currency     string             `json:"currency"`
	Charge       *string            `json:"charge"`
	Order        primitive.ObjectID `json:"order"`
}

type ThirdPartyBalanceResponse struct {
	Status   string `json:"status"`
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
}
