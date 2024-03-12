package user

import (
	"smm/internal/api/serializers"
	"smm/internal/database/models"
	"smm/internal/database/queries"
	"smm/pkg/bcrypt"
	"smm/pkg/constants"
	"smm/pkg/request"
	"smm/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller interface {
	Create(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	GenerateApiKey(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	UpdateBalance(ctx *fiber.Ctx) error
	UpdatePassword(ctx *fiber.Ctx) error
}
type controller struct{}

func New() Controller {
	return &controller{}
}

func (ctrl *controller) Create(ctx *fiber.Ctx) error {
	var requestBody serializers.UserCreateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	user, err := userQuery.CreateOne(models.User{
		FirstName:         requestBody.FirstName,
		LastName:          requestBody.LastName,
		Username:          requestBody.Username,
		PhoneNumber:       requestBody.PhoneNumber,
		Language:          requestBody.Language,
		Status:            requestBody.Status,
		Balance:           requestBody.Balance,
		Email:             requestBody.Email,
		EmailVerification: false,
		Address:           requestBody.Address,
		Password:          bcrypt.GeneratePassword(requestBody.Password),
		Avatar:            requestBody.Avatar,
		ApiKey:            utils.UUIDv4(),
	})
	if err != nil {
		return err
	}
	return response.New(ctx, response.Response{StatusCode: fiber.StatusCreated, Data: fiber.Map{"id": user.Id}})
}

func (ctrl *controller) List(ctx *fiber.Ctx) error {
	var requestBody serializers.UserListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField(
		"updated_at", "created_at", "last_active", "username",
		"phone_number", "language", "status", "balance", "email", "email_verification",
		"2fa_enable", "address", "avatar", "api_key", "_id")
	totalChan := make(chan int64, 1)
	errChan := make(chan error, 1)
	go func() {
		total, err := userQuery.GetTotalByFilter(requestBody.GetFilter())
		if err != nil {
			errChan <- err
			return
		}
		totalChan <- total
		errChan <- nil
	}()
	queryOption.AddSort(requestBody.Sort())
	users, err := userQuery.GetByFilter(requestBody.GetFilter(), queryOption)
	if err != nil {
		return err
	}
	if err = <-errChan; err != nil {
		return err
	}
	res := make([]serializers.UserListResponse, len(users))
	for i, user := range users {
		res[i].CreatedAt = user.CreatedAt
		res[i].UpdatedAt = user.UpdatedAt
		res[i].Status = user.Status
		res[i].LastActive = user.LastActive
		res[i].Username = user.Username
		res[i].PhoneNumber = user.PhoneNumber
		res[i].Language = user.Language
		res[i].Status = user.Status
		res[i].Balance = user.Balance
		res[i].Email = user.Email
		res[i].EmailVerification = user.EmailVerification
		res[i].TwoFAEnable = user.TwoFAEnable
		res[i].Address = user.Address
		res[i].Avatar = user.Avatar
		res[i].ApiKey = user.ApiKey
		res[i].Id = user.Id

	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}

func (ctrl *controller) GenerateApiKey(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	userQuery := queries.NewUser(ctx.Context())
	apiKey := utils.UUIDv4()
	if err := userQuery.UpdateApiKeyById(userId, apiKey); err != nil {
		return err
	}
	return response.New(ctx, response.Response{StatusCode: fiber.StatusOK, Data: fiber.Map{"api_key": apiKey}})
}

func (ctrl *controller) Update(ctx *fiber.Ctx) error {
	var requestBody serializers.UserUpdateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	if err := userQuery.UpdateById(requestBody.Id, queries.UserUpdateByIdDoc{
		FirstName:         requestBody.FirstName,
		LastName:          requestBody.LastName,
		Username:          requestBody.Username,
		PhoneNumber:       requestBody.PhoneNumber,
		Language:          requestBody.Language,
		Status:            requestBody.Status,
		Email:             requestBody.Email,
		EmailVerification: requestBody.EmailVerification,
		TwoFAEnable:       requestBody.TwoFAEnable,
		Address:           requestBody.Address,
		Avatar:            requestBody.Avatar,
	}); err != nil {
		return err
	}
	return response.New(ctx, response.Response{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	userQuery := queries.NewUser(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField(
		"updated_at", "created_at", "last_active", "first_name", "last_name", "username",
		"phone_number", "language", "status", "balance", "email", "email_verification",
		"2fa_enable", "address", "avatar", "api_key", "_id")
	user, err := userQuery.GetById(userId, queryOption)
	if err != nil {
		return err
	}
	return response.New(ctx, response.Response{StatusCode: fiber.StatusOK, Data: serializers.UserGetResponse{
		UpdatedAt:         user.UpdatedAt,
		CreatedAt:         user.CreatedAt,
		LastActive:        user.LastActive,
		Username:          user.Username,
		PhoneNumber:       user.PhoneNumber,
		Language:          user.Language,
		Status:            user.Status,
		Balance:           user.Balance,
		Email:             user.Email,
		EmailVerification: user.EmailVerification,
		TwoFAEnable:       user.TwoFAEnable,
		Address:           user.Address,
		Avatar:            user.Avatar,
		ApiKey:            user.ApiKey,
		Id:                user.Id,
	}})
}

func (ctrl *controller) UpdateBalance(ctx *fiber.Ctx) error {
	var requestBody serializers.UserUpdateBalanceBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	if err := userQuery.UpdateBalanceById(requestBody.Id, requestBody.Balance); err != nil {
		return err
	}
	return response.New(ctx, response.Response{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) UpdatePassword(ctx *fiber.Ctx) error {
	var requestBody serializers.UserUpdatePasswordBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	if err := userQuery.UpdatePasswordById(requestBody.Id, bcrypt.GeneratePassword(requestBody.Password)); err != nil {
		return err
	}
	return response.New(ctx, response.Response{StatusCode: fiber.StatusOK})
}
