package category

import (
	"smm/internal/api/serializers"
	"smm/internal/database/models"
	"smm/internal/database/queries"
	"smm/pkg/constants"
	"smm/pkg/request"
	"smm/pkg/response"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller interface {
	Create(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	UserList(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}
type controller struct{}

func New() Controller {
	return &controller{}
}

func (ctrl *controller) Create(ctx *fiber.Ctx) error {
	var requestBody serializers.CategoryCreateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	categoryQuery := queries.NewCategory(ctx.Context())
	category, err := categoryQuery.CreateOne(models.Category{
		Title:       requestBody.Title,
		Description: requestBody.Description,
		Image:       requestBody.Image,
		Status:      constants.CategoryStatusOn,
	})
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusCreated, Data: fiber.Map{"id": category.Id}})
}

func (ctrl *controller) List(ctx *fiber.Ctx) error {
	var requestBody serializers.CategoryListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
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
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}

func (ctrl *controller) Update(ctx *fiber.Ctx) error {
	var requestBody serializers.CategoryUpdateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	categoryQuery := queries.NewCategory(ctx.Context())
	if err := categoryQuery.UpdateById(requestBody.Id, queries.CategoryUpdateByIdDoc{
		Title:       requestBody.Title,
		Image:       requestBody.Image,
		Description: requestBody.Description,
		Status:      requestBody.Status,
	}); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	categoryId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	categoryQuery := queries.NewCategory(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField("title", "description", "created_at", "updated_at", "_id", "status", "image")
	category, err := categoryQuery.GetById(categoryId, queryOption)
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: serializers.CategoryGetResponse{
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
		Title:       category.Title,
		Description: category.Description,
		Id:          category.Id,
		Status:      category.Status,
		Image:       category.Image,
	}})
}

func (ctrl *controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	categoryId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	categoryQuery := queries.NewCategory(ctx.Context())
	if err := categoryQuery.DeleteById(categoryId); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) UserList(ctx *fiber.Ctx) error {
	var requestBody serializers.CategoryUserListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	categoryQuery := queries.NewCategory(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField("title", "_id")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)
	filter := requestBody.GetFilter()
	filter["status"] = constants.CategoryStatusOn

	go func() {
		total, err := categoryQuery.GetTotalByFilter(filter)
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	categories, err := categoryQuery.GetByFilter(filter, queryOption)
	if err != nil {
		return err
	}
	if err = <-errChan; err != nil {
		return err
	}
	res := make([]serializers.CategoryUserListResponse, len(categories))
	for i, category := range categories {
		res[i].Title = category.Title
		res[i].Id = category.Id

	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}
