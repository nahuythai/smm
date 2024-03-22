package user

import (
	"smm/internal/api/serializers"
	"smm/internal/database/models"
	"smm/internal/database/queries"
	"smm/pkg/bcrypt"
	"smm/pkg/configure"
	"smm/pkg/constants"
	"smm/pkg/jwt"
	"smm/pkg/logging"
	"smm/pkg/mail"
	"smm/pkg/otp"
	"smm/pkg/request"
	"smm/pkg/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	cfg    = configure.GetConfig()
	logger = logging.GetLogger()
)

type Controller interface {
	Create(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	GenerateApiKey(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Me(ctx *fiber.Ctx) error
	UpdateBalance(ctx *fiber.Ctx) error
	UpdatePassword(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	VerifyLogin(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
	VerifyEmail(ctx *fiber.Ctx) error
}
type controller struct {
	s *service
}

func New() Controller {
	return &controller{
		s: &service{},
	}
}

func (ctrl *controller) Create(ctx *fiber.Ctx) error {
	var requestBody serializers.UserCreateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	user, err := userQuery.CreateOne(models.User{
		DisplayName:       requestBody.DisplayName,
		Username:          requestBody.Username,
		PhoneNumber:       requestBody.PhoneNumber,
		Language:          requestBody.Language,
		Status:            requestBody.Status,
		Balance:           requestBody.Balance,
		Email:             requestBody.Email,
		EmailVerification: false,
		Address:           requestBody.Address,
		Password:          bcrypt.GeneratePassword(requestBody.Password),
		TwoFAEnable:       requestBody.TwoFAEnable,
		Avatar:            requestBody.Avatar,
		ApiKey:            utils.UUIDv4(),
		Role:              constants.UserRoleUser,
	})
	if err != nil {
		return err
	}
	secretKey, _ := otp.GetGlobal().GenerateSecretKey(requestBody.Username)
	if _, err = queries.NewOtp(ctx.Context()).CreateOne(models.Otp{
		UserId:    user.Id,
		SecretKey: secretKey.Secret(),
		OtpUrl:    secretKey.URL(),
	}); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusCreated, Data: fiber.Map{"id": user.Id}})
}

func (ctrl *controller) List(ctx *fiber.Ctx) error {
	var requestBody serializers.UserListBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	queryOption := queries.NewOption()
	pagination := request.NewPagination(requestBody.Page, requestBody.Limit)
	queryOption.SetPagination(pagination)
	queryOption.SetOnlyField(
		"username", "display_name", "phone_number", "status", "balance", "email", "avatar", "_id")
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
		res[i].Username = user.Username
		res[i].PhoneNumber = user.PhoneNumber
		res[i].Status = user.Status
		res[i].Balance = user.Balance
		res[i].Email = user.Email
		res[i].Avatar = user.Avatar
		res[i].Id = user.Id
		res[i].DisplayName = user.DisplayName

	}
	pagination.SetTotal(<-totalChan)
	return response.NewPaginationResponse(ctx, response.PaginationResponse{StatusCode: fiber.StatusOK, Data: res, Extras: *pagination})
}

func (ctrl *controller) GenerateApiKey(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	userQuery := queries.NewUser(ctx.Context())
	apiKey := utils.UUIDv4()
	if err := userQuery.UpdateApiKeyById(userId, apiKey); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: fiber.Map{"api_key": apiKey}})
}

func (ctrl *controller) Update(ctx *fiber.Ctx) error {
	var requestBody serializers.UserUpdateBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	if err := userQuery.UpdateById(requestBody.Id, queries.UserUpdateByIdDoc{
		DisplayName:       requestBody.DisplayName,
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
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) Get(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	userQuery := queries.NewUser(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField(
		"updated_at", "created_at", "last_active", "display_name", "username",
		"phone_number", "language", "status", "balance", "email", "email_verification",
		"2fa_enable", "address", "avatar", "api_key", "_id")
	user, err := userQuery.GetById(userId, queryOption)
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: serializers.UserGetResponse{
		DisplayName:       user.DisplayName,
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
		Role:              user.Role,
	}})
}

func (ctrl *controller) UpdateBalance(ctx *fiber.Ctx) error {
	var requestBody serializers.UserUpdateBalanceBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	backgroundQuery := queries.NewBackground(ctx.Context())
	task, err := backgroundQuery.CreateOne(models.Background{
		Type:   constants.BackgroundTypeUpdateBalance,
		UserId: requestBody.Id,
	})
	if err != nil {
		return err
	}
	if err := userQuery.UpdateBalanceById(requestBody.Id, requestBody.Balance); err != nil {
		return err
	}
	if err = backgroundQuery.DeleteById(task.Id); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) UpdatePassword(ctx *fiber.Ctx) error {
	var requestBody serializers.UserUpdatePasswordBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	if err := userQuery.UpdatePasswordById(requestBody.Id, bcrypt.GeneratePassword(requestBody.Password)); err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK})
}

func (ctrl *controller) Login(ctx *fiber.Ctx) error {
	var requestBody serializers.UserLoginBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField("password", "2fa_enable", "_id", "email_verification", "email")
	user, err := userQuery.GetByUsername(requestBody.Username, queryOption)
	if err != nil {
		if err.(*response.Option).Code == constants.ErrCodeUserNotFound {
			return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeAppUnauthorized})
		}
		return err
	}
	if !bcrypt.ComparePassword(user.Password, requestBody.Password) {
		return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeAppUnauthorized})
	}
	if !user.EmailVerification {
		session, err := queries.NewSession(ctx.Context()).CreateOne(models.Session{
			ExpiredAt: time.Now().Add(cfg.VerifyEmailSessionDuration),
			UserId:    user.Id,
			Type:      constants.SessionTypeVerifyEmail,
		})
		if err != nil {
			return err
		}
		token, err := jwt.GetGlobal().GenerateToken(session.Id, constants.TokenTypeSession, cfg.VerifyEmailSessionDuration)
		if err != nil {
			return err
		}
		// send mail
		go func() {
			mailBody, err := ctrl.s.createVerifyEmailTemplate(requestBody.Username, token)
			if err != nil {
				logger.Error().Err(err).Caller().Str("func", "Login").Str("funcInline", "ctrl.s.CreateVerifyEmailTemplate").Msg("user-controller")
			}
			if err = mail.GetGlobal().SendHTMLMail(user.Email, mailBody); err != nil {
				logger.Error().Err(err).Caller().Str("func", "Login").Str("funcInline", "mail.New().SendHTMLMail").Msg("user-controller")
			}
		}()
		return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeUserInvalidEmail})
	}
	if user.TwoFAEnable {
		session, err := queries.NewSession(ctx.Context()).CreateOne(models.Session{
			ExpiredAt: time.Now().Add(cfg.SessionDuration),
			UserId:    user.Id,
			Type:      constants.SessionTypeVerifyLogin,
		})
		if err != nil {
			return err
		}
		token, err := jwt.GetGlobal().GenerateToken(session.Id, constants.TokenTypeSession, cfg.SessionDuration)
		if err != nil {
			return err
		}
		return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: fiber.Map{"token": token, "token_type": constants.TokenTypeSession}})
	}
	token, err := jwt.GetGlobal().GenerateToken(user.Id, constants.TokenTypeAccess, cfg.AccessTokenDuration)
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: fiber.Map{"token": token, "token_type": constants.TokenTypeAccess}})
}

func (ctrl *controller) VerifyLogin(ctx *fiber.Ctx) error {
	session := ctx.Locals(constants.LocalSessionKey).(*models.Session)
	if session.Type != constants.SessionTypeVerifyLogin {
		return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeTokenWrong, Data: "missing or wrong token"})
	}
	var requestBody serializers.UserLoginVerifyBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	otpCode, err := queries.NewOtp(ctx.Context()).GetByUserId(session.UserId)
	if err != nil {
		return err
	}
	if !otp.GetGlobal().Validate(requestBody.Code, otpCode.SecretKey) {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgOtpWrong, Code: constants.ErrCodeOtpWrong})
	}
	if err := queries.NewSession(ctx.Context()).DeleteById(session.Id); err != nil {
		return err
	}
	token, err := jwt.GetGlobal().GenerateToken(session.Id, constants.TokenTypeAccess, cfg.SessionDuration)
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: fiber.Map{"token": token, "token_type": constants.TokenTypeAccess}})
}

func (ctrl *controller) Register(ctx *fiber.Ctx) error {
	var requestBody serializers.UserRegisterBodyValidate
	if err := ctx.BodyParser(&requestBody); err != nil {
		return response.NewError(fiber.StatusBadRequest, response.Option{Data: constants.ErrMsgFieldWrongType, Code: constants.ErrCodeAppBadRequest})
	}
	if err := requestBody.Validate(); err != nil {
		return err
	}
	userQuery := queries.NewUser(ctx.Context())
	user, err := userQuery.CreateOne(models.User{
		DisplayName:       requestBody.DisplayName,
		Username:          requestBody.Username,
		PhoneNumber:       requestBody.PhoneNumber,
		Status:            constants.UserStatusBanned,
		Balance:           0,
		Email:             requestBody.Email,
		EmailVerification: false,
		Password:          bcrypt.GeneratePassword(requestBody.Password),
		TwoFAEnable:       false,
		ApiKey:            utils.UUIDv4(),
	})
	if err != nil {
		return err
	}
	return response.New(ctx, response.Option{StatusCode: fiber.StatusCreated, Data: fiber.Map{"id": user.Id}})
}

func (ctrl *controller) Me(ctx *fiber.Ctx) error {
	user := ctx.Locals(constants.LocalUserKey).(*models.User)
	return response.New(ctx, response.Option{StatusCode: fiber.StatusOK, Data: serializers.UserGetResponse{
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
		Role:              user.Role,
	}})
}

func (ctrl *controller) VerifyEmail(ctx *fiber.Ctx) error {
	ctx.Context().SetContentType(fiber.MIMETextHTML)
	session, err := ctrl.s.sessionValidate(ctx.Context(), ctx.Query("token"), constants.SessionTypeVerifyEmail)
	if err != nil {
		body, _ := ctrl.s.verifyEmailFailTemplate()
		return ctx.SendString(body)
	}
	if err := queries.NewUser(ctx.Context()).UpdateEmailVerificationById(session.UserId, true); err != nil {
		body, _ := ctrl.s.verifyEmailFailTemplate()
		return ctx.SendString(body)
	}
	if err := queries.NewSession(ctx.Context()).DeleteById(session.Id); err != nil {
		body, _ := ctrl.s.verifyEmailFailTemplate()
		return ctx.SendString(body)
	}
	body, _ := ctrl.s.verifyEmailSuccessTemplate()
	return ctx.SendString(body)
}
