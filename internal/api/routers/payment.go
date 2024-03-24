package routers

import (
	paymentCtrl "smm/internal/api/controller/payment"
	"smm/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

type Payment interface {
	V1()
}

type payment struct {
	ctrl   paymentCtrl.Controller
	router fiber.Router
}

func NewPayment(router fiber.Router) Payment {
	return &payment{
		ctrl:   paymentCtrl.New(),
		router: router,
	}
}

func (r *payment) V1() {
	r.User()
	r.Admin()
}

func (r *payment) User() {
	router := r.router.Group("/payments")
	router.Use(middleware.UserAuth)
	router.Post("/qr-top-up", r.ctrl.QRTopUp)
	router.Post("/list", r.ctrl.ListByUser)
}

func (r *payment) Admin() {
	router := r.router.Group("/admin/payments")
	router.Use(middleware.UserAuth, middleware.AdminPermission)
	router.Get("/:id", r.ctrl.Get)
	router.Put("", r.ctrl.Update)
	router.Post("/list", r.ctrl.List)
}
