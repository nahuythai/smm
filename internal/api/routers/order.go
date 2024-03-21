package routers

import (
	orderCtrl "smm/internal/api/controller/order"
	"smm/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

type Order interface {
	V1()
}

type order struct {
	ctrl   orderCtrl.Controller
	router fiber.Router
}

func NewOrder(router fiber.Router) Order {
	return &order{
		ctrl:   orderCtrl.New(),
		router: router,
	}
}

func (r *order) V1() {
	r.User()
	r.Admin()
}

func (r *order) User() {
	router := r.router.Group("/orders")
	router.Use(middleware.UserAuth)
	router.Post("/", r.ctrl.Create)
	router.Post("/list", r.ctrl.ListByUser)
}

func (r *order) Admin() {
	router := r.router.Group("admin/orders")
	router.Use(middleware.UserAuth, middleware.AdminPermission)
	router.Post("/list", r.ctrl.List)
	router.Put("/", r.ctrl.Update)
	router.Get("/:id", r.ctrl.Get)
	router.Delete("/:id", r.ctrl.Delete)
}
