package paymentmethod

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
	var requestBody serializers.PaymentMethodCreateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	paymentMethodQuery := queries.NewPaymentMethod(ctx.Context())
	paymentMethod, err := paymentMethodQuery.CreateOne(models.PaymentMethod{
		Status:           constants.PaymentMethodStatusOn,
		Name:             requestBody.Name,
		Code:             requestBody.Code,
		Image:            requestBody.Image,
		MinAmount:        requestBody.MinAmount,
		MaxAmount:        requestBody.MaxAmount,
		PercentageCharge: requestBody.PercentageCharge,
		Description:      requestBody.Description,
		Auto:             *requestBody.Auto,
		AccountName:      requestBody.AccountName,
		AccountNumber:    requestBody.AccountNumber,
	})
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusCreated, Data: fiber.Map{"id": paymentMethod.Id}})
}

func (ctrl *controller) List(ctx *fiber.Ctx) error {
	var requestBody serializers.PaymentMethodListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	paymentMethodQuery := queries.NewPaymentMethod(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField("name", "code", "image", "status", "_id")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)

	go func() {
		total, err := paymentMethodQuery.GetTotalByFilter(requestBody.GetFilter())
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	paymentMethods, err := paymentMethodQuery.GetByFilter(requestBody.GetFilter(), queryOption)
	if err != nil {
		return err
	}
	if err = <-errChan; err != nil {
		return err
	}
	res := make([]serializers.PaymentMethodListResponse, len(paymentMethods))
	for i, paymentMethod := range paymentMethods {
		res[i].Name = paymentMethod.Name
		res[i].Id = paymentMethod.Id
		res[i].Status = paymentMethod.Status
		res[i].Code = paymentMethod.Code
		res[i].Image = paymentMethod.Image
	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}

func (ctrl *controller) Update(ctx *fiber.Ctx) error {
	var requestBody serializers.PaymentMethodUpdateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	paymentMethodQuery := queries.NewPaymentMethod(ctx.Context())
	if err := paymentMethodQuery.UpdateById(requestBody.Id, queries.PaymentMethodUpdateByIdDoc{
		Name:             requestBody.Name,
		Code:             requestBody.Code,
		Image:            requestBody.Image,
		MinAmount:        requestBody.MinAmount,
		MaxAmount:        requestBody.MaxAmount,
		PercentageCharge: requestBody.PercentageCharge,
		Status:           requestBody.Status,
		Description:      requestBody.Description,
		Auto:             *requestBody.Auto,
		AccountName:      requestBody.AccountName,
		AccountNumber:    requestBody.AccountNumber,
	}); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	paymentMethodId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	paymentMethodQuery := queries.NewPaymentMethod(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField("_id", "name", "code", "min_amount", "max_amount", "auto",
		"percentage_charge", "description", "status", "image",
		"fixed_charge", "convention_rate", "account_name", "account_number")
	paymentMethod, err := paymentMethodQuery.GetById(paymentMethodId, queryOption)
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: serializers.PaymentMethodGetResponse{
		Name:             paymentMethod.Name,
		Code:             paymentMethod.Code,
		Image:            paymentMethod.Image,
		MinAmount:        paymentMethod.MinAmount,
		MaxAmount:        paymentMethod.MaxAmount,
		PercentageCharge: paymentMethod.PercentageCharge,
		FixedCharge:      paymentMethod.FixedCharge,
		ConventionRate:   paymentMethod.ConventionRate,
		Description:      paymentMethod.Description,
		Auto:             paymentMethod.Auto,
		Status:           paymentMethod.Status,
		Id:               paymentMethod.Id,
		AccountName:      paymentMethod.AccountName,
		AccountNumber:    paymentMethod.AccountNumber,
	}})
}

func (ctrl *controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	paymentMethodId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	paymentMethodQuery := queries.NewPaymentMethod(ctx.Context())
	if err := paymentMethodQuery.DeleteById(paymentMethodId); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}
