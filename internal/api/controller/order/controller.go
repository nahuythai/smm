package order

import (
	"context"
	"fmt"
	"smm/internal/api/serializers"
	"smm/internal/database/models"
	"smm/internal/database/queries"
	"smm/pkg/constants"
	providerapi "smm/pkg/providerapi"
	"smm/pkg/request"
	"smm/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller interface {
	Create(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	ListByUser(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}
type controller struct{}

func New() Controller {
	return &controller{}
}

func (ctrl *controller) Create(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals(constants.LocalUserKey).(*models.User)
	var requestBody serializers.OrderCreateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}

	queryOption := queries.NewOption()
	queryOption.SetOnlyField("provider_service_id", "provider_id", "category_id", "rate", "_id", "min_amount", "max_amount")
	service, err := queries.NewService(ctx.Context()).GetActiveById(requestBody.ServiceId, queryOption)
	if err != nil {
		return err
	}

	if requestBody.Quantity < service.MinAmount || requestBody.Quantity > service.MaxAmount {
		return response.NewError(fiber.StatusBadRequest, response.Option{Code: constants.ErrCodeOrderQuantityInvalid})
	}
	queryOption.SetOnlyField("url", "api_key")
	provider, err := queries.NewProvider(ctx.Context()).GetActiveById(service.ProviderId, queryOption)
	if err != nil {
		return err
	}

	// TODO: get custom rate and rate limit
	customRateQuery := queries.NewCustomRate(ctx.Context())
	queryOption.SetOnlyField("price")
	customRate, err := customRateQuery.GetByUserIdAndServiceId(currentUser.Id, service.Id, queryOption)
	if err != nil {
		if err.(*response.Option).Code != constants.ErrCodeCustomRateNotFound {
			return err
		}
	}
	serviceRate := service.Rate
	if customRate != nil {
		serviceRate = customRate.Price
	}
	backgroundQuery := queries.NewBackground(ctx.Context())
	task, err := backgroundQuery.CreateOne(models.Background{
		Type:   constants.BackgroundTypeUpdateBalance,
		UserId: currentUser.Id,
	})
	if err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	queryOption.SetOnlyField("balance")
	user, err := userQuery.GetById(currentUser.Id, queryOption)
	if err != nil {
		return err
	}
	price := float64(requestBody.Quantity) * serviceRate / 1000
	if user.Balance < price {
		return response.NewError(fiber.StatusBadRequest, response.Option{Code: constants.ErrCodeUserNotEnoughBalance, Data: "not enough balance"})
	}
	if err = userQuery.UpdateBalanceById(currentUser.Id, currentUser.Balance-price); err != nil {
		return err
	}
	if err = backgroundQuery.DeleteById(task.Id); err != nil {
		return err
	}

	order := models.Order{
		Quantity:   requestBody.Quantity,
		Status:     constants.OrderStatusProcessing,
		Link:       requestBody.Link,
		UserId:     currentUser.Id,
		ServiceId:  requestBody.ServiceId,
		Price:      price,
		CategoryId: service.CategoryId,
		ProviderId: service.ProviderId,
	}

	orderQuery := queries.NewOrder(ctx.Context())
	newOrder, err := orderQuery.CreateOne(order)
	if err != nil {
		return err
	}
	if _, err := queries.NewTransaction(ctx.Context()).CreateOne(models.Transaction{
		Type:   constants.TransactionTypePlayOrder,
		UserId: currentUser.Id,
		Amount: order.Price,
	}); err != nil {
		return err
	}
	go func() {
		res, _ := providerapi.New(provider.Url, provider.ApiKey).AddOrder(providerapi.AddOrderRequest{
			ApiServiceId: service.ProviderServiceId,
			Link:         requestBody.Link,
			Quantity:     requestBody.Quantity,
		})
		if res != nil {
			if res.Error == "" {
				_ = queries.NewOrder(context.Background()).UpdateProviderOrderIdAndProviderResponseById(newOrder.Id, res.OrderId, fmt.Sprintf("Order: %d", res.OrderId))
			} else {
				_ = queries.NewOrder(context.Background()).UpdateProviderOrderIdAndProviderResponseById(newOrder.Id, res.OrderId, fmt.Sprintf("Error: %s", res.Error))
			}
		}
	}()
	return response.New(ctx, response.Option{StatusCode: fiber.StatusCreated, Data: fiber.Map{"id": newOrder.Id}})
}

func (ctrl *controller) List(ctx *fiber.Ctx) error {
	var requestBody serializers.OrderListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	orderQuery := queries.NewOrder(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField("_id", "link", "updated_at", "quantity",
		"status", "start_counter", "remains", "price", "service_id", "user_id")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)

	go func() {
		total, err := orderQuery.GetTotalByFilter(requestBody.GetFilter())
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	orders, err := orderQuery.GetByFilter(requestBody.GetFilter(), queryOption)
	if err != nil {
		return err
	}
	serviceIds := lo.Map(orders, func(order models.Order, i int) primitive.ObjectID {
		return order.ServiceId
	})
	serviceIds = lo.Uniq(serviceIds)
	if err = <-errChan; err != nil {
		return err
	}

	queryOption.SetOnlyField("title")
	services, err := queries.NewService(ctx.Context()).GetByIds(serviceIds, queryOption)
	if err != nil {
		return err
	}
	serviceNames := lo.SliceToMap(services, func(service models.Service) (primitive.ObjectID, string) {
		return service.Id, service.Title
	})
	userIds := lo.Map(orders, func(order models.Order, _ int) primitive.ObjectID {
		return order.UserId
	})
	queryOption.SetOnlyField("username")
	users, err := queries.NewUser(ctx.Context()).GetByIds(userIds, queryOption)
	if err != nil {
		return err
	}
	usernames := lo.SliceToMap(users, func(user models.User) (primitive.ObjectID, string) {
		return user.Id, user.Username
	})

	res := make([]serializers.OrderListResponse, len(orders))
	for i, order := range orders {
		res[i].Link = order.Link
		res[i].UpdatedAt = order.UpdatedAt
		res[i].Id = order.Id
		res[i].Status = order.Status
		res[i].Price = order.Price
		res[i].Quantity = order.Quantity
		res[i].Remains = order.Remains
		res[i].StartCounter = order.StartCounter
		res[i].Service = serviceNames[order.ServiceId]
		res[i].Username = usernames[order.UserId]

	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}

func (ctrl *controller) ListByUser(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals(constants.LocalUserKey).(*models.User)
	var requestBody serializers.OrderListByUserBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	orderQuery := queries.NewOrder(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField("_id", "link", "created_at", "quantity",
		"status", "start_counter", "remains", "price", "service_id")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)

	filter := requestBody.GetFilter()
	filter["user_id"] = currentUser.Id
	go func() {
		total, err := orderQuery.GetTotalByFilter(filter)
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	orders, err := orderQuery.GetByFilter(filter, queryOption)
	if err != nil {
		return err
	}
	serviceIds := lo.Map(orders, func(order models.Order, i int) primitive.ObjectID {
		return order.ServiceId
	})
	serviceIds = lo.Uniq(serviceIds)
	if err = <-errChan; err != nil {
		return err
	}

	queryOption.SetOnlyField("title")
	services, err := queries.NewService(ctx.Context()).GetByIds(serviceIds, queryOption)
	if err != nil {
		return err
	}
	serviceNames := lo.SliceToMap(services, func(service models.Service) (primitive.ObjectID, string) {
		return service.Id, service.Title
	})

	res := make([]serializers.OrderListByUserResponse, len(orders))
	for i, order := range orders {
		res[i].Link = order.Link
		res[i].CreatedAt = order.CreatedAt
		res[i].Id = order.Id
		res[i].Status = order.Status
		res[i].Price = order.Price
		res[i].Quantity = order.Quantity
		res[i].Remains = order.Remains
		res[i].StartCounter = order.StartCounter
		res[i].Service = serviceNames[order.ServiceId]

	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}

func (ctrl *controller) Update(ctx *fiber.Ctx) error {
	var requestBody serializers.OrderUpdateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	orderQuery := queries.NewOrder(ctx.Context())
	if err := orderQuery.UpdateById(requestBody.Id, queries.OrderUpdateByIdDoc{
		Status:       requestBody.Status,
		StartCounter: requestBody.StartCounter,
		Remains:      requestBody.Remains,
		Note:         requestBody.Note,
		Link:         requestBody.Link,
	}); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	orderId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	orderQuery := queries.NewOrder(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField("quantity", "provider_order_id", "provider_response", "status", "_id", "start_counter",
		"remains", "note", "link", "service_id")
	order, err := orderQuery.GetById(orderId, queryOption)
	if err != nil {
		return err
	}
	queryOption.SetOnlyField("title", "category_id", "provider_id")
	service, err := queries.NewService(ctx.Context()).GetById(order.ServiceId, queryOption)
	if err != nil {
		return err
	}
	queryOption.SetOnlyField("title")
	categoryTitle := ""
	category, _ := queries.NewCategory(ctx.Context()).GetById(service.CategoryId, queryOption)
	if category != nil {
		categoryTitle = category.Title
	}
	queryOption.SetOnlyField("api_name")
	provider, err := queries.NewProvider(ctx.Context()).GetById(service.ProviderId, queryOption)
	if err != nil {
		return err
	}

	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: serializers.OrderGetResponse{
		Quantity:              order.Quantity,
		ProviderOrderId:       order.ProviderOrderId,
		ProviderOrderResponse: order.ProviderOrderResponse,
		Status:                order.Status,
		StartCounter:          order.StartCounter,
		Remains:               order.Remains,
		Note:                  order.Note,
		Link:                  order.Link,
		Service:               service.Title,
		Category:              categoryTitle,
		Provider:              provider.ApiName,
		Id:                    order.Id,
	}})
}

func (ctrl *controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	orderId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	orderQuery := queries.NewOrder(ctx.Context())
	if err := orderQuery.DeleteById(orderId); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}
