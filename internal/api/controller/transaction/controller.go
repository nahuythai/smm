package transaction

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
	List(ctx *fiber.Ctx) error
	UserList(ctx *fiber.Ctx) error
}
type controller struct{}

func New() Controller {
	return &controller{}
}

func (ctrl *controller) List(ctx *fiber.Ctx) error {
	var requestBody serializers.TransactionListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	transactionQuery := queries.NewTransaction(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField("created_at", "trx_id", "amount", "type", "user_id")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)
	filter := requestBody.GetFilter()
	go func() {
		total, err := transactionQuery.GetTotalByFilter(filter)
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	transactions, err := transactionQuery.GetByFilter(filter, queryOption)
	if err != nil {
		return err
	}
	if err = <-errChan; err != nil {
		return err
	}
	userIds := lo.Map(transactions, func(transaction models.Transaction, _ int) primitive.ObjectID {
		return transaction.UserId
	})
	queryOption.SetOnlyField("username")
	users, err := queries.NewUser(ctx.Context()).GetByIds(userIds, queryOption)
	if err != nil {
		return err
	}
	usernames := lo.SliceToMap(users, func(user models.User) (primitive.ObjectID, string) {
		return user.Id, user.Username
	})
	res := make([]serializers.TransactionListResponse, len(transactions))
	for i, transaction := range transactions {
		res[i].CreatedAt = transaction.CreatedAt
		res[i].TransactionId = transaction.TransactionId
		res[i].Amount = transaction.Amount
		res[i].Type = transaction.Type
		res[i].Username = usernames[transaction.UserId]
	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}

func (ctrl *controller) UserList(ctx *fiber.Ctx) error {
	currentUser := ctx.Locals(constants.LocalUserKey).(*models.User)
	var requestBody serializers.TransactionListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	transactionQuery := queries.NewTransaction(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField("created_at", "trx_id", "amount", "type")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)
	filter := requestBody.GetFilter()
	filter["user_id"] = currentUser.Id
	go func() {
		total, err := transactionQuery.GetTotalByFilter(filter)
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	transactions, err := transactionQuery.GetByFilter(filter, queryOption)
	if err != nil {
		return err
	}
	if err = <-errChan; err != nil {
		return err
	}
	res := make([]serializers.TransactionUserListResponse, len(transactions))
	for i, transaction := range transactions {
		res[i].CreatedAt = transaction.CreatedAt
		res[i].TransactionId = transaction.TransactionId
		res[i].Amount = transaction.Amount
		res[i].Type = transaction.Type
	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}
