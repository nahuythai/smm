package webhook

import (
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	Web2MMBBank(ctx *fiber.Ctx) error
}
type controller struct{}

func New() Controller {
	return &controller{}
}

func (ctrl *controller) Web2MMBBank(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"status": true})
}
