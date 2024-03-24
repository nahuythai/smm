package routers

import (
	webhookCtrl "smm/internal/api/controller/webhook"

	"github.com/gofiber/fiber/v2"
)

type Webhook interface {
	V1()
}

type webhook struct {
	ctrl   webhookCtrl.Controller
	router fiber.Router
}

func NewWebhook(router fiber.Router) Webhook {
	return &webhook{
		ctrl:   webhookCtrl.New(),
		router: router,
	}
}

func (r *webhook) V1() {
	r.Web2M()
}

func (r *webhook) Web2M() {
	router := r.router.Group("/webhooks/web2m")
	// router.Use(middleware.UserAuth)
	router.Post("/mb-bank", r.ctrl.Web2MMBBank)
}
