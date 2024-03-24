package routers

import (
	categoryCtrl "smm/internal/api/controller/category"
	"smm/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

type Category interface {
	V1()
}

type category struct {
	ctrl   categoryCtrl.Controller
	router fiber.Router
}

func NewCategory(router fiber.Router) Category {
	return &category{
		ctrl:   categoryCtrl.New(),
		router: router,
	}
}

func (r *category) V1() {
	r.Admin()
	r.User()
}

func (r *category) Admin() {
	router := r.router.Group("/admin/categories")
	router.Use(middleware.UserAuth, middleware.AdminPermission)
	router.Post("/", r.ctrl.Create)
	router.Post("/list", r.ctrl.List)
	router.Put("/", r.ctrl.Update)
	router.Get("/:id", r.ctrl.Get)
	router.Delete("/:id", r.ctrl.Delete)
}

func (r *category) User() {
	router := r.router.Group("/categories")
	router.Use(middleware.UserAuth)
	router.Post("/list", r.ctrl.UserList)
}
