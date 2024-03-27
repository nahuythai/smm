package thirdparty

import (
	"context"
	"fmt"
	"smm/internal/api/serializers"
	"smm/internal/database/models"
	"smm/internal/database/queries"
	"smm/pkg/constants"
	providerapi "smm/pkg/providerapi"
	"smm/pkg/response"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller interface {
	Route(ctx *fiber.Ctx) error
}
type controller struct{}

func New() Controller {
	return &controller{}
}

func (ctrl *controller) Route(ctx *fiber.Ctx) error {
	var requestBody serializers.ThirdPartyRouteValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "body wrong format"})
	}
	if err := requestBody.Validate(); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	switch requestBody.Action {
	case constants.ThirdPartyActionOrderCreate:
		return ctrl.CreateOrder(ctx)
	case constants.ThirdPartyActionOrderStatus:
		return ctrl.OrderStatus(ctx)
	case constants.ThirdPartyActionOrderMultipleStatus:
		return ctrl.MultipleOrderStatus(ctx)
	case constants.ThirdPartyActionUserBalance:
		return ctrl.Balance(ctx)
	case constants.ThirdPartyActionServiceList:
		return ctrl.ListService(ctx)
	default:
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "the selected action is invalid"})
	}
}

func (ctrl *controller) Balance(ctx *fiber.Ctx) error {
	user := ctx.Locals(constants.LocalUserKey).(*models.User)
	return ctx.JSON(serializers.ThirdPartyBalanceResponse{
		Status:   "success",
		Balance:  fmt.Sprintf("%v", user.Balance),
		Currency: "VND",
	})
}

func (ctrl *controller) CreateOrder(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals(constants.LocalUserKey).(*models.User)
	var requestBody serializers.ThirdPartyOrderCreateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "body wrong format"})
	}
	if err := requestBody.Validate(); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	queryOption := queries.NewOption()
	queryOption.SetOnlyField("provider_service_id", "provider_id", "category_id", "rate", "_id", "min_amount", "max_amount")
	service, err := queries.NewService(ctx.Context()).GetActiveBySeq(requestBody.Service, queryOption)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	if requestBody.Quantity < service.MinAmount || requestBody.Quantity > service.MaxAmount {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "quantity invalid"})
	}
	queryOption.SetOnlyField("url", "api_key")
	provider, err := queries.NewProvider(ctx.Context()).GetActiveById(service.ProviderId, queryOption)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	// TODO: get custom rate and rate limit
	customRateQuery := queries.NewCustomRate(ctx.Context())
	queryOption.SetOnlyField("price")
	customRate, err := customRateQuery.GetByUserIdAndServiceId(currentUser.Id, service.Id, queryOption)
	if err != nil {
		if err.(*response.Option).Code != constants.ErrCodeCustomRateNotFound {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
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
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	userQuery := queries.NewUser(ctx.Context())
	queryOption.SetOnlyField("balance")
	user, err := userQuery.GetById(currentUser.Id, queryOption)
	if err != nil {
		return err
	}
	price := float64(requestBody.Quantity) * serviceRate / 1000
	if user.Balance < price {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "not enough balance"})
	}
	if err = userQuery.UpdateBalanceById(currentUser.Id, currentUser.Balance-price); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	if err = backgroundQuery.DeleteById(task.Id); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	order := models.Order{
		Quantity:   requestBody.Quantity,
		Status:     constants.OrderStatusProcessing,
		Link:       requestBody.Link,
		UserId:     currentUser.Id,
		ServiceId:  service.Id,
		Price:      price,
		CategoryId: service.CategoryId,
	}

	orderQuery := queries.NewOrder(ctx.Context())
	newOrder, err := orderQuery.CreateOne(order)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	if _, err := queries.NewTransaction(ctx.Context()).CreateOne(models.Transaction{
		Type:   constants.TransactionTypePlayOrder,
		UserId: currentUser.Id,
		Amount: order.Price,
	}); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
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
	return response.New(ctx, response.Option{StatusCode: fiber.StatusCreated, Data: fiber.Map{"status": "success", "order": newOrder.Id}})
}

func (ctrl *controller) OrderStatus(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals(constants.LocalUserKey).(*models.User)
	var requestBody serializers.ThirdPartyOrderStatusBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "body wrong format"})
	}
	if err := requestBody.Validate(); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	orderId, err := primitive.ObjectIDFromHex(requestBody.OrderId)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "order not found"})
	}
	orderQuery := queries.NewOrder(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField("quantity", "status", "_id", "start_counter", "remains")
	order, err := orderQuery.GetByIdAndUserId(orderId, currentUser.Id, queryOption)
	if err != nil {
		return err
	}
	return ctx.JSON(serializers.ThirdPartyOrderStatusResponse{
		Quantity:     order.Quantity,
		Status:       constants.OrderStatusTextMapping[order.Status],
		StartCounter: order.StartCounter,
		Remains:      order.Remains,
		Currency:     "VND",
		Charge:       nil,
	})
}

func (ctrl *controller) MultipleOrderStatus(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals(constants.LocalUserKey).(*models.User)
	var requestBody serializers.ThirdPartyMultipleOrderStatusBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "body wrong format"})
	}
	if err := requestBody.Validate(); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	orderStringIds := strings.Split(requestBody.Orders, ",")
	orderIds := make([]primitive.ObjectID, 0, len(orderStringIds))
	for _, orderId := range orderStringIds {
		if value, err := primitive.ObjectIDFromHex(orderId); err == nil {
			orderIds = append(orderIds, value)
		}
	}
	orderQuery := queries.NewOrder(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField("quantity", "status", "_id", "start_counter", "remains")
	orders, err := orderQuery.GetByIdsAndUserId(orderIds, currentUser.Id, queryOption)
	if err != nil {
		return err
	}
	result := make([]serializers.ThirdPartyMultipleOrderStatusResponse, 0, len(orders))
	for _, order := range orders {
		result = append(result, serializers.ThirdPartyMultipleOrderStatusResponse{
			Quantity:     order.Quantity,
			Status:       constants.OrderStatusTextMapping[order.Status],
			StartCounter: order.StartCounter,
			Remains:      order.Remains,
			Currency:     "VND",
			Charge:       nil,
			Order:        order.Id,
		})
	}
	return ctx.JSON(result)
}

func (ctrl *controller) ListService(ctx *fiber.Ctx) error {
	var requestBody serializers.ServiceUserListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	serviceQuery := queries.NewService(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField("title", "seq", "_id", "category_id", "min_amount", "max_amount", "rate")
	filter := requestBody.GetFilter()
	filter["status"] = constants.ServiceStatusOn
	queryOption.AddSort(requestBody.Sort())
	services, err := serviceQuery.GetByFilter(bson.M{"status": constants.ServiceStatusOn}, queryOption)
	if err != nil {
		return err
	}
	categoryIds := make([]primitive.ObjectID, 0, len(services))
	for _, service := range services {
		categoryIds = append(categoryIds, service.CategoryId)
	}
	queryOption.SetOnlyField("title", "_id")
	categoryIdNameMapping := make(map[primitive.ObjectID]string)
	categories, err := queries.NewCategory(ctx.Context()).GetByIds(categoryIds)
	if err != nil {
		return err
	}
	for _, category := range categories {
		categoryIdNameMapping[category.Id] = category.Title
	}
	res := make([]serializers.ThirdPartListServiceResponse, len(services))
	for i, service := range services {
		res[i].Title = service.Title
		res[i].Seq = service.Seq
		res[i].Category = categoryIdNameMapping[service.CategoryId]
		res[i].MaxAmount = service.MaxAmount
		res[i].MinAmount = service.MinAmount
		res[i].Rate = fmt.Sprintf("%v", service.Rate)
	}
	return ctx.JSON(res)
}
