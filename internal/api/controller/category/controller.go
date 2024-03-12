package category

import (
	"smm/internal/api/serializers"
	"smm/internal/database/models"
	"smm/internal/database/queries"
	"smm/pkg/constants"
	"smm/pkg/request"
	"smm/pkg/response"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type Controller interface {
	Create(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
}
type controller struct{}

func New() Controller {
	return &controller{}
}

func (ctrl *controller) Create(ctx *fiber.Ctx) error {
	var requestBody serializers.CategoryCreateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: "Field wrong type", Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	categoryQuery := queries.NewCategory(ctx.Context())
	category, err := categoryQuery.CreateOne(models.Category{
		Title:       requestBody.Title,
		Description: requestBody.Description,
		Image:       requestBody.Image,
		Status:      requestBody.Status,
	})
	if err != nil {
		return err
	}
	return response.New(ctx, response.Response{StatusCode: fiber.StatusCreated, Data: fiber.Map{"id": category.Id}})
}

func (ctrl *controller) List(ctx *fiber.Ctx) error {
	var requestBody serializers.CategoryListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: "Field wrong type", Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	categoryQuery := queries.NewCategory(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField("title", "description", "created_at", "updated_at", "_id", "status", "image")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)

	go func() {
		total, err := categoryQuery.GetTotalByFilter(requestBody.GetFilter())
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	categories, err := categoryQuery.GetByFilter(requestBody.GetFilter(), queryOption)
	if err != nil {
		return err
	}
	if err = <-errChan; err != nil {
		return err
	}
	res := make([]serializers.CategoryListResponse, len(categories))
	for i, category := range categories {
		res[i].CreatedAt = category.CreatedAt
		res[i].UpdatedAt = category.UpdatedAt
		res[i].Title = category.Title
		res[i].Description = category.Description
		res[i].Id = category.Id
		res[i].Status = category.Status
		res[i].Image = category.Image

	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusCreated, Data: res, Extras: *pagination})
}

func (ctrl *controller) Update(ctx *fiber.Ctx) error {
	var requestBody serializers.CategoryUpdateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: "Field wrong type", Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	categoryQuery := queries.NewCategory(ctx.Context())
	if err := categoryQuery.UpdateById(requestBody.Id,
		bson.M{
			"title":       requestBody.Title,
			"description": requestBody.Description,
			"image":       requestBody.Image,
			"status":      requestBody.Status,
		}); err != nil {
		return err
	}
	return response.New(ctx, response.Response{StatusCode: fiber.StatusOK})
}
