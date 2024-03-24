package service

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
	var requestBody serializers.ServiceCreateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	serviceQuery := queries.NewService(ctx.Context())
	service, err := serviceQuery.CreateOne(models.Service{
		Title:             requestBody.Title,
		Status:            constants.ServiceStatusOn,
		MinAmount:         requestBody.MinAmount,
		MaxAmount:         requestBody.MaxAmount,
		Rate:              requestBody.Rate,
		Description:       requestBody.Description,
		ProviderServiceId: requestBody.ProviderServiceId,
		CategoryId:        requestBody.CategoryId,
		ProviderId:        requestBody.ProviderId,
	})
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusCreated, Data: fiber.Map{"id": service.Id}})
}

func (ctrl *controller) List(ctx *fiber.Ctx) error {
	var requestBody serializers.ServiceListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	serviceQuery := queries.NewService(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField("title", "description", "_id", "status", "provider_id", "min_amount", "max_amount", "rate")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)

	go func() {
		total, err := serviceQuery.GetTotalByFilter(requestBody.GetFilter())
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	services, err := serviceQuery.GetByFilter(requestBody.GetFilter(), queryOption)
	if err != nil {
		return err
	}
	if err = <-errChan; err != nil {
		return err
	}
	providerIds := make([]primitive.ObjectID, 0, len(services))
	for _, service := range services {
		providerIds = append(providerIds, service.ProviderId)
	}
	queryOption.SetOnlyField("api_name", "_id")
	providerIdNameMapping := make(map[primitive.ObjectID]string)
	providers, err := queries.NewProvider(ctx.Context()).GetByIds(providerIds)
	if err != nil {
		return err
	}
	for _, provider := range providers {
		providerIdNameMapping[provider.Id] = provider.ApiName
	}
	res := make([]serializers.ServiceListResponse, len(services))
	for i, service := range services {
		res[i].Title = service.Title
		res[i].Id = service.Id
		res[i].Status = service.Status
		res[i].Provider = providerIdNameMapping[service.ProviderId]
		res[i].MaxAmount = service.MaxAmount
		res[i].MinAmount = service.MinAmount
		res[i].Description = service.Description
	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}

func (ctrl *controller) Update(ctx *fiber.Ctx) error {
	var requestBody serializers.ServiceUpdateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	queryOption := queries.NewOption()
	queryOption.SetOnlyField("_id")
	if _, err := queries.NewCategory(ctx.Context()).GetById(requestBody.CategoryId, queryOption); err != nil {
		return err
	}
	if _, err := queries.NewProvider(ctx.Context()).GetById(requestBody.ProviderId, queryOption); err != nil {
		return err
	}

	serviceQuery := queries.NewService(ctx.Context())
	if err := serviceQuery.UpdateById(requestBody.Id, queries.ServiceUpdateByIdDoc{
		Title:             requestBody.Title,
		Status:            requestBody.Status,
		MinAmount:         requestBody.MinAmount,
		MaxAmount:         requestBody.MaxAmount,
		Rate:              requestBody.Rate,
		Description:       requestBody.Description,
		ProviderServiceId: requestBody.ProviderServiceId,
		CategoryId:        requestBody.CategoryId,
		ProviderId:        requestBody.ProviderId,
	}); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	serviceId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	serviceQuery := queries.NewService(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField("title", "status", "min_amount", "max_amount", "rate", "description", "provider_service_id", "category_id", "provider_id")
	service, err := serviceQuery.GetById(serviceId, queryOption)
	if err != nil {
		return err
	}
	queryOption.SetOnlyField("title")
	categoryName := ""
	category, err := queries.NewCategory(ctx.Context()).GetById(service.CategoryId, queryOption)
	if err != nil {
		if err.(*response.Option).Code != constants.ErrCodeCategoryNotFound {
			return err
		}
	}
	if category != nil {
		categoryName = category.Title
	}

	providerName := ""
	queryOption.SetOnlyField("api_name")
	provider, err := queries.NewProvider(ctx.Context()).GetById(service.ProviderId, queryOption)
	if err != nil {
		if err.(*response.Option).Code != constants.ErrCodeProviderNotFound {
			return err
		}
	}
	if provider != nil {
		providerName = provider.ApiName
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: serializers.ServiceGetResponse{
		Title:             service.Title,
		Status:            service.Status,
		MinAmount:         service.MinAmount,
		MaxAmount:         service.MaxAmount,
		Rate:              service.Rate,
		Description:       service.Description,
		ProviderServiceId: service.ProviderServiceId,
		Category:          categoryName,
		Provider:          providerName,
		Id:                service.Id,
	}})
}

func (ctrl *controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	serviceId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	serviceQuery := queries.NewService(ctx.Context())
	if err := serviceQuery.DeleteById(serviceId); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) UserList(ctx *fiber.Ctx) error {
	var requestBody serializers.ServiceUserListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	serviceQuery := queries.NewService(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField("title", "description", "_id", "status", "provider_id", "min_amount", "max_amount", "rate")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)
	filter := requestBody.GetFilter()
	filter["status"] = constants.ServiceStatusOn

	go func() {
		total, err := serviceQuery.GetTotalByFilter(filter)
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	services, err := serviceQuery.GetByFilter(filter, queryOption)
	if err != nil {
		return err
	}
	if err = <-errChan; err != nil {
		return err
	}
	providerIds := make([]primitive.ObjectID, 0, len(services))
	for _, service := range services {
		providerIds = append(providerIds, service.ProviderId)
	}
	queryOption.SetOnlyField("api_name", "_id")
	providerIdNameMapping := make(map[primitive.ObjectID]string)
	providers, err := queries.NewProvider(ctx.Context()).GetByIds(providerIds)
	if err != nil {
		return err
	}
	for _, provider := range providers {
		providerIdNameMapping[provider.Id] = provider.ApiName
	}
	res := make([]serializers.ServiceUserListResponse, len(services))
	for i, service := range services {
		res[i].Title = service.Title
		res[i].Id = service.Id
		res[i].Status = service.Status
		res[i].Provider = providerIdNameMapping[service.ProviderId]
		res[i].MaxAmount = service.MaxAmount
		res[i].MinAmount = service.MinAmount
		res[i].Description = service.Description
	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}
