package payment

import (
	"smm/internal/api/serializers"
	"smm/internal/database/models"
	"smm/internal/database/queries"
	"smm/pkg/constants"
	"smm/pkg/request"
	"smm/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller interface {
	QRTopUp(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	ListByUser(ctx *fiber.Ctx) error
}
type controller struct{}

func New() Controller {
	return &controller{}
}

func (ctrl *controller) QRTopUp(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals(constants.LocalUserKey).(*models.User)
	var requestBody serializers.PaymentQRTopUpBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	queryOption := queries.NewOption()
	queryOption.SetOnlyField("account_number", "account_name", "code", "min_amount", "max_amount")
	paymentMethod, err := queries.NewPaymentMethod(ctx.Context()).GetById(requestBody.Method, queryOption)
	if err != nil {
		return err
	}
	if requestBody.Amount < paymentMethod.MinAmount || requestBody.Amount > paymentMethod.MaxAmount {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrCodePaymentAmountInvalid})
	}
	paymentQuery := queries.NewPayment(ctx.Context())
	payment, err := paymentQuery.CreateOne(models.Payment{
		Status:        constants.PaymentStatusPending,
		Amount:        requestBody.Amount,
		UserId:        currentUser.Id,
		Method:        requestBody.Method,
		TransactionId: new(models.Transaction).GenerateTransactionId(),
	})
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusCreated, Data: serializers.PaymentQRTopUpResponse{
		AccountName:   paymentMethod.AccountName,
		AccountNumber: paymentMethod.AccountNumber,
		TransactionId: payment.TransactionId,
		Amount:        payment.Amount,
		MethodCode:    paymentMethod.Code,
	}})
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
		Status: constants.PaymentStatusPending,
		Amount: requestBody.Amount,
		UserId: requestBody.UserId,
		Method: requestBody.Method,
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
	queryOption.SetOnlyField("amount", "trx_id", "method", "status", "_id", "feedback", "user_id")
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
	paymentMethodIds := lo.Map(payments, func(payment models.Payment, i int) primitive.ObjectID {
		return payment.Method
	})
	paymentMethodIds = lo.Uniq(paymentMethodIds)

	queryOption.SetOnlyField("name")
	methods, err := queries.NewPaymentMethod(ctx.Context()).GetByIds(paymentMethodIds, queryOption)
	if err != nil {
		return err
	}
	methodNames := lo.SliceToMap(methods, func(method models.PaymentMethod) (primitive.ObjectID, string) {
		return method.Id, method.Name
	})

	userIds := lo.Map(payments, func(payment models.Payment, i int) primitive.ObjectID {
		return payment.UserId
	})
	userIds = lo.Uniq(userIds)
	queryOption.SetOnlyField("username")
	users, err := queries.NewUser(ctx.Context()).GetByIds(userIds, queryOption)
	if err != nil {
		return err
	}
	usernames := lo.SliceToMap(users, func(user models.User) (primitive.ObjectID, string) {
		return user.Id, user.Username
	})
	if err = <-errChan; err != nil {
		return err
	}
	res := make([]serializers.PaymentListResponse, len(payments))
	for i, payment := range payments {
		res[i].Id = payment.Id
		res[i].Status = payment.Status
		res[i].Method = methodNames[payment.Method]
		res[i].Amount = payment.Amount
		res[i].Feedback = payment.Feedback
		res[i].TransactionId = payment.TransactionId
		res[i].Username = usernames[payment.UserId]
	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}

func (ctrl *controller) ListByUser(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals(constants.LocalUserKey).(*models.User)
	var requestBody serializers.PaymentListByUserBodyValidate
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
	queryOption.SetOnlyField("amount", "trx_id", "method", "status", "_id")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)
	filter := requestBody.GetFilter()
	filter["user_id"] = currentUser.Id

	go func() {
		total, err := paymentQuery.GetTotalByFilter(filter)
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	payments, err := paymentQuery.GetByFilter(filter, queryOption)
	if err != nil {
		return err
	}
	paymentMethodIds := lo.Map(payments, func(payment models.Payment, i int) primitive.ObjectID {
		return payment.Method
	})
	paymentMethodIds = lo.Uniq(paymentMethodIds)

	queryOption.SetOnlyField("name")
	methods, err := queries.NewPaymentMethod(ctx.Context()).GetByIds(paymentMethodIds, queryOption)
	if err != nil {
		return err
	}
	methodNames := lo.SliceToMap(methods, func(method models.PaymentMethod) (primitive.ObjectID, string) {
		return method.Id, method.Name
	})
	if err = <-errChan; err != nil {
		return err
	}
	res := make([]serializers.PaymentListByUserResponse, len(payments))
	for i, payment := range payments {
		res[i].Id = payment.Id
		res[i].Status = payment.Status
		res[i].Method = methodNames[payment.Method]
		res[i].Amount = payment.Amount
		res[i].TransactionId = payment.TransactionId
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
	queryOption.SetOnlyField("_id", "created_at", "amount", "status", "feedback", "trx_id", "user_id", "method")
	payment, err := paymentQuery.GetById(paymentId, queryOption)
	if err != nil {
		return err
	}
	queryOption.SetOnlyField("username")
	user, err := queries.NewUser(ctx.Context()).GetById(payment.UserId, queryOption)
	if err != nil {
		return err
	}

	queryOption.SetOnlyField("name")
	method, err := queries.NewPaymentMethod(ctx.Context()).GetById(payment.Method, queryOption)
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: serializers.PaymentGetResponse{
		CreatedAt:     payment.CreatedAt,
		Amount:        payment.Amount,
		Status:        payment.Status,
		Feedback:      payment.Feedback,
		TransactionId: payment.TransactionId,
		Username:      user.Username,
		Method:        method.Name,
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
