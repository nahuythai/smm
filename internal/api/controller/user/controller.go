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
)

type Controller interface {
	Create(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
}
type controller struct{}

func New() Controller {
	return &controller{}
}

func (ctrl *controller) Create(ctx *fiber.Ctx) error {
	var requestBody serializers.UserCreateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: "Field wrong type", Code: constants.ErrCodeAppBadRequest})
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
		ApiToken:          utils.UUIDv4(),
	})
	if err != nil {
		return err
	}
	return response.New(ctx, response.Response{StatusCode: fiber.StatusCreated, Data: fiber.Map{"id": user.Id}})
}

func (ctrl *controller) List(ctx *fiber.Ctx) error {
	var requestBody serializers.UserListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.ErrorResponse{Err: "Field wrong type", Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField(
		"updated_at", "created_at", "last_active", "first_name", "last_name", "username",
		"phone_number", "language", "status", "balance", "email", "email_verification",
		"2fa_enable", "address", "avatar", "api_token", "_id")
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
		res[i].FirstName = user.FirstName
		res[i].LastName = user.LastName
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
		res[i].ApiToken = user.ApiToken
		res[i].Id = user.Id

	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusCreated, Data: res, Extras: *pagination})
}
