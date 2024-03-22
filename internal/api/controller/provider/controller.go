package provider

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
	Update(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}
type controller struct{}

func New() Controller {
	return &controller{}
}

func (ctrl *controller) Create(ctx *fiber.Ctx) error {
	var requestBody serializers.ProviderCreateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	providerQuery := queries.NewProvider(ctx.Context())
	provider, err := providerQuery.CreateOne(models.Provider{
		Status:      constants.ProviderStatusOn,
		ApiName:     requestBody.ApiName,
		ApiKey:      requestBody.ApiKey,
		Description: requestBody.Description,
		Rate:        requestBody.Rate,
		Url:         requestBody.Url,
		Balance:     0,
	})
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusCreated, Data: fiber.Map{"id": provider.Id}})
}

func (ctrl *controller) List(ctx *fiber.Ctx) error {
	var requestBody serializers.ProviderListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	providerQuery := queries.NewProvider(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField("api_name", "description", "_id", "status", "balance")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)

	go func() {
		total, err := providerQuery.GetTotalByFilter(requestBody.GetFilter())
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	categories, err := providerQuery.GetByFilter(requestBody.GetFilter(), queryOption)
	if err != nil {
		return err
	}
	if err = <-errChan; err != nil {
		return err
	}
	res := make([]serializers.ProviderListResponse, len(categories))
	for i, provider := range categories {
		res[i].ApiName = provider.ApiName
		res[i].Status = provider.Status
		res[i].Balance = provider.Balance
		res[i].Description = provider.Description
		res[i].Id = provider.Id

	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}

func (ctrl *controller) Update(ctx *fiber.Ctx) error {
	var requestBody serializers.ProviderUpdateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	providerQuery := queries.NewProvider(ctx.Context())
	if err := providerQuery.UpdateById(requestBody.Id, queries.ProviderUpdateByIdDoc{
		ApiName:     requestBody.ApiName,
		ApiKey:      requestBody.ApiKey,
		Description: requestBody.Description,
		Status:      *requestBody.Status,
		Rate:        requestBody.Rate,
		Url:         requestBody.Url,
	}); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	providerId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	providerQuery := queries.NewProvider(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField("api_name", "description", "_id", "api_key", "url", "status", "balance", "rate")
	provider, err := providerQuery.GetById(providerId, queryOption)
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: serializers.ProviderGetResponse{
		Status:      provider.Status,
		ApiName:     provider.ApiName,
		ApiKey:      provider.ApiKey,
		Description: provider.Description,
		Url:         provider.Url,
		Balance:     provider.Balance,
		Rate:        provider.Rate,
		Id:          provider.Id,
	}})
}

func (ctrl *controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	providerId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	providerQuery := queries.NewProvider(ctx.Context())
	if err := providerQuery.DeleteById(providerId); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}
