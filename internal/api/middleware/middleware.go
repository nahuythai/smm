package middleware

import (
	"smm/internal/database/models"
	"smm/internal/database/queries"
	"smm/pkg/constants"
	"smm/pkg/jwt"
	"smm/pkg/logging"
	"smm/pkg/response"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	logger = logging.GetLogger()
)

func SessionAuth(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")
	token, _ = strings.CutPrefix(token, "Bearer ")
	payload, err := jwt.GetGlobal().ValidateToken(token)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "SessionAuth").Str("funcInline", "jwt.GetGlobal().ValidateToken").Msg("session-middleware")
		return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeTokenWrong, Data: "missing or wrong token"})
	}
	if payload.Type != constants.TokenTypeSession {
		return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeTokenWrong, Data: "missing or wrong token"})
	}
	id, _ := primitive.ObjectIDFromHex(payload.ID)
	session, err := queries.NewSession(ctx.Context()).GetById(id)
	if err != nil {
		if e, ok := err.(*response.Option); ok {
			if e.Code == constants.ErrCodeSessionNotFound {
				return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeSessionNotFound})
			}
		}
	}
	ctx.Locals(constants.LocalSessionKey, session)
	return ctx.Next()
}

func UserAuth(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")
	token, _ = strings.CutPrefix(token, "Bearer ")
	payload, err := jwt.GetGlobal().ValidateToken(token)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "UserAuth").Str("funcInline", "jwt.GetGlobal().ValidateToken").Msg("user-middleware")
		return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeTokenWrong, Data: "missing or wrong token"})
	}
	if payload.Type != constants.TokenTypeAccess {
		return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeTokenWrong, Data: "missing or wrong token"})
	}
	id, _ := primitive.ObjectIDFromHex(payload.ID)
	userQuery := queries.NewUser(ctx.Context())
	queryOption := queries.NewOption()
	queryOption.SetOnlyField(
		"updated_at", "created_at", "last_active", "first_name", "last_name", "username",
		"phone_number", "language", "status", "balance", "email", "email_verification",
		"2fa_enable", "address", "avatar", "api_key", "_id", "role")
	user, err := userQuery.GetById(id, queryOption)
	if err != nil {
		if e, ok := err.(*response.Option); ok {
			if e.Code == constants.ErrCodeUserNotFound {
				return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeUserNotFound})
			}
		}
	}
	if user.Status == constants.UserStatusBanned {
		return response.NewError(fiber.StatusForbidden, response.Option{Code: constants.ErrCodeUserBanned})
	}
	ctx.Locals(constants.LocalUserKey, user)
	return ctx.Next()
}

func AdminPermission(ctx *fiber.Ctx) error {
	user := ctx.Locals(constants.LocalUserKey).(*models.User)
	if user.Role != constants.UserRoleAdmin {
		return response.NewError(fiber.StatusForbidden, response.Option{Code: constants.ErrCodeUserAdminPermissionRequired})
	}
	return ctx.Next()
}
