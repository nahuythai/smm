package payment

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

func (ctrl *controller) QRTopUp(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals(constants.LocalUserKey).(*models.User)
	var requestBody serializers.PaymentTopUpBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	// transaction, err := queries.NewTransaction(ctx.Context()).CreateOne(models.Transaction{
	// 	UserId: currentUser.Id,
	// 	Amount: requestBody.Amount,
	// 	TransactionId: utils.UUID(),
	// })
	// if err != nil {
	// 	return err
	// }
	paymentQuery := queries.NewPayment(ctx.Context())
	payment, err := paymentQuery.CreateOne(models.Payment{
		Status:        constants.PaymentStatusPending,
		Amount:        requestBody.Amount,
		UserId:        currentUser.Id,
		PaymentMethod: requestBody.PaymentMethod,
	})
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusCreated, Data: fiber.Map{"id": payment.Id}})
}

func (ctrl *controller) Create(ctx *fiber.Ctx) error {
	var requestBody serializers.PaymentCreateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	paymentQuery := queries.NewPayment(ctx.Context())
	payment, err := paymentQuery.CreateOne(models.Payment{
		Status:        constants.PaymentStatusPending,
		Amount:        requestBody.Amount,
		UserId:        requestBody.UserId,
		PaymentMethod: requestBody.PaymentMethod,
	})
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusCreated, Data: fiber.Map{"id": payment.Id}})
}

func (ctrl *controller) List(ctx *fiber.Ctx) error {
	var requestBody serializers.PaymentListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	paymentQuery := queries.NewPayment(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField("amount", "trx_id", "payment_method", "status", "_id", "feedback", "user_id")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)

	go func() {
		total, err := paymentQuery.GetTotalByFilter(requestBody.GetFilter())
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	payments, err := paymentQuery.GetByFilter(requestBody.GetFilter(), queryOption)
	if err != nil {
		return err
	}
	if err = <-errChan; err != nil {
		return err
	}
	res := make([]serializers.PaymentListResponse, len(payments))
	for i, payment := range payments {
		res[i].Id = payment.Id
		res[i].Status = payment.Status
		res[i].PaymentMethod = nil
		res[i].Amount = payment.Amount
		res[i].Feedback = payment.Feedback
		res[i].TransactionId = payment.TransactionId
		res[i].Username = ""
	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}

func (ctrl *controller) Update(ctx *fiber.Ctx) error {
	var requestBody serializers.PaymentUpdateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	paymentQuery := queries.NewPayment(ctx.Context())
	if err := paymentQuery.UpdateById(requestBody.Id, queries.PaymentUpdateByIdDoc{
		Status:   requestBody.Status,
		Feedback: requestBody.Feedback,
	}); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	paymentId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	paymentQuery := queries.NewPayment(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField("_id", "created_at", "amount", "status", "feedback", "trx_id", "user_id")
	payment, err := paymentQuery.GetById(paymentId, queryOption)
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: serializers.PaymentGetResponse{
		CreatedAt:     payment.CreatedAt,
		Amount:        payment.Amount,
		Status:        payment.Status,
		Feedback:      payment.Feedback,
		TransactionId: payment.TransactionId,
		Username:      "",
		PaymentMethod: nil,
		Id:            payment.Id,
	}})
}

func (ctrl *controller) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	paymentId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	paymentQuery := queries.NewPayment(ctx.Context())
	if err := paymentQuery.DeleteById(paymentId); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}
