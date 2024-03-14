package middleware

import (
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

func TransactionAuth(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")
	token, _ = strings.CutPrefix(token, "Bearer ")
	payload, err := jwt.GetGlobal().ValidateToken(token)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "TransactionAuth").Str("funcInline", "jwt.GetGlobal().ValidateToken").Msg("transaction-middleware")
		return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeTokenWrong, Data: "missing or wrong token"})
	}
	if payload.Type != constants.TokenTypeTransaction {
		return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeTokenWrong, Data: "missing or wrong token"})
	}
	id, _ := primitive.ObjectIDFromHex(payload.ID)
	transaction, err := queries.NewTransaction(ctx.Context()).GetById(id)
	if err != nil {
		if e, ok := err.(*response.Option); ok {
			if e.Code == constants.ErrCodeTransactionNotFound {
				return response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeTransactionNotFound})
			}
		}
	}
	ctx.Locals(constants.LocalTransactionKey, transaction)
	return ctx.Next()
}
