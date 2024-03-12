package serializers

import (
	"fmt"
	"regexp"
	"smm/internal/database/queries"
	"smm/pkg/constants"
	"smm/pkg/response"
	"smm/pkg/validator"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryCreateBodyValidate struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"omitempty"`
	Image       string `json:"image" validate:"omitempty"`
	Status      int    `json:"status" validate:"omitempty"`
}

func (v *CategoryCreateBodyValidate) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}
	if v.Status != constants.CategoryStatusOn && v.Status != constants.CategoryStatusOff {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: "status invalid", Code: constants.ErrCodeAppBadRequest})
	}
	return nil
}

type CategoryListBodyValidate struct {
	Page     int64                `json:"page" validate:"omitempty"`
	Limit    int64                `json:"limit" validate:"omitempty"`
	SortType int                  `json:"sort_type" validate:"omitempty,oneof=-1 1 0"`
	SortBy   string               `json:"sort_by" validate:"omitempty"`
	Search   string               `json:"search" validate:"omitempty"`
	Filters  []CategoryListFilter `json:"filters" validate:"omitempty,dive"`
}

func (v *CategoryListBodyValidate) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}
	return nil
}

func (v *CategoryListBodyValidate) Sort() map[string]int {
	sortBy := "updated_at"
	if v.SortBy != "" {
		sortBy = v.SortBy
	}
	return map[string]int{sortBy: v.SortType}
}

type CategoryListFilter struct {
	Value  interface{} `json:"value" validate:"required"`
	Method string      `json:"method" validate:"required,oneof=$eq $regex"`
	Field  string      `json:"field" validate:"required,oneof=username email phone_number"`
}

func (v *CategoryListBodyValidate) GetFilter() bson.M {
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

type CategoryListResponse struct {
	UpdatedAt   time.Time          `json:"updated_at"`
	CreatedAt   time.Time          `json:"created_at"`
	Title       string             `json:"title"`
	Image       string             `json:"image"`
	Status      int                `json:"status"`
	Description string             `json:"description"`
	Id          primitive.ObjectID `json:"id"`
}

type CategoryUpdateBodyValidate struct {
	Id          primitive.ObjectID `json:"id" validate:"required"`
	Title       string             `json:"title" validate:"required"`
	Description *string            `json:"description" validate:"omitempty"`
	Image       *string            `json:"image" validate:"omitempty"`
	Status      *int               `json:"status" validate:"omitempty,oneof=0 1"`
}

func (v *CategoryUpdateBodyValidate) Validate() error {
	if err := validator.Validate(v); err != nil {
		return err
	}
	return nil
}

type CategoryGetResponse struct {
	UpdatedAt   time.Time          `json:"updated_at"`
	CreatedAt   time.Time          `json:"created_at"`
	Title       string             `json:"title"`
	Image       string             `json:"image"`
	Status      int                `json:"status"`
	Description string             `json:"description"`
	Id          primitive.ObjectID `json:"id"`
}
