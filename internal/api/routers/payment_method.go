package routers

import (
	paymentMethodCtrl "smm/internal/api/controller/paymentmethod"
	"smm/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

type PaymentMethod interface {
	V1()
}

type paymentMethod struct {
	ctrl   paymentMethodCtrl.Controller
	router fiber.Router
}

func NewPaymentMethod(router fiber.Router) PaymentMethod {
	return &paymentMethod{
		ctrl:   paymentMethodCtrl.New(),
		router: router,
	}
}

func (r *paymentMethod) V1() {
	r.Admin()
}

func (r *paymentMethod) Admin() {
	router := r.router.Group("/admin/payment-methods")
	router.Use(middleware.UserAuth, middleware.AdminPermission)
	router.Post("/", r.ctrl.Create)
	router.Post("/list", r.ctrl.List)
	router.Put("/", r.ctrl.Update)
	router.Get("/:id", r.ctrl.Get)
	router.Delete("/:id", r.ctrl.Delete)
}
