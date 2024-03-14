package response

import (
	"errors"
	"fmt"
	"smm/pkg/constants"
	"smm/pkg/request"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type Option struct {
	StatusCode int
	Data       interface{}
	Code       string
}
type Response struct {
	StatusCode int
	Data       interface{}
}

type ErrorResponse struct {
	StatusCode int
	Code       string
	Err        string
}

func (e *Option) Error() string {
	return fmt.Sprintf("%v", e.Data)
}

func NewError(statusCode int, opts ...Option) error {
	if len(opts) == 0 {
		return &Option{
			StatusCode: statusCode,
			Data:       utils.StatusMessage(statusCode),
			Code:       constants.ErrCodeAppUnknown,
		}
	}
	if opts[0].Code == "" {
		opts[0].Code = constants.ErrCodeAppUnknown
	}
	opts[0].StatusCode = statusCode
	return &opts[0]
}

func New(ctx *fiber.Ctx, opt Option) error {
	data := fiber.Map{"success": true}
	if opt.Data != nil {
		data["data"] = opt.Data
	}
	return ctx.Status(opt.StatusCode).JSON(data)
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
	var e *Option
	if errors.As(err, &e) {
		if e.Data == nil {
			e.Data = utils.StatusMessage(e.StatusCode)
		}
		return ctx.Status(e.StatusCode).JSON(fiber.Map{"code": e.Code, "error": e.Error(), "success": false})
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
