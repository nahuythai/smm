package response

import (
	"errors"
	"smm/pkg/constants"
	"smm/pkg/request"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	StatusCode int
	Data       map[string]interface{}
}

type ErrorResponse struct {
	StatusCode int
	Code       string
	Err        string
}

func (e *ErrorResponse) Error() string {
	return e.Err
}

func NewError(statusCode int, res ...ErrorResponse) error {
	if len(res) == 0 {
		return &ErrorResponse{
			StatusCode: statusCode,
		}
	}
	res[0].StatusCode = statusCode
	return &res[0]
}

func New(ctx *fiber.Ctx, res Response) error {
	data := fiber.Map{"success": true}
	if res.Data != nil {
		data["data"] = res.Data
	}
	return ctx.Status(res.StatusCode).JSON(data)
}

type PaginationResponse struct {
	StatusCode int
	Data       interface{}
	Extras     request.Pagination
}

func NewPaginationResponse(ctx *fiber.Ctx, res PaginationResponse) error {
	data := fiber.Map{"success": true, "data": make([]struct{}, 0)}
	if res.Data != nil {
		data["data"] = res.Data
	}
	data["extras"] = res.Extras
	return ctx.Status(res.StatusCode).JSON(data)
}

func FiberErrorHandler(ctx *fiber.Ctx, err error) error {
	var e *ErrorResponse
	if errors.As(err, &e) {
		return ctx.Status(e.StatusCode).JSON(fiber.Map{"code": e.Code, "error": e.Err, "success": false})
	} else if fiberErr, ok := err.(*fiber.Error); ok {
		return ctx.Status(fiberErr.Code).JSON(fiber.Map{"error": fiberErr.Message, "success": false, "code": constants.ErrCodeAppUnknown})

	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(
		fiber.Map{
			"code":    constants.ErrCodeAppInternalServerError,
			"error":   fiber.ErrInternalServerError.Message,
			"success": false,
		})
}
