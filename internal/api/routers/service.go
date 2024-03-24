package routers

import (
	serviceCtrl "smm/internal/api/controller/service"
	"smm/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

type Service interface {
	V1()
}

type service struct {
	ctrl   serviceCtrl.Controller
	router fiber.Router
}

func NewService(router fiber.Router) Service {
	return &service{
		ctrl:   serviceCtrl.New(),
		router: router,
	}
}

func (r *service) V1() {
	r.Admin()
	r.User()
}

func (r *service) Admin() {
	router := r.router.Group("/admin/services")
	router.Use(middleware.UserAuth, middleware.AdminPermission)
	router.Post("/", r.ctrl.Create)
	router.Post("/list", r.ctrl.List)
	router.Put("/", r.ctrl.Update)
	router.Get("/:id", r.ctrl.Get)
	router.Delete("/:id", r.ctrl.Delete)
}

func (r *service) User() {
	router := r.router.Group("/services")
	router.Use(middleware.UserAuth)
	router.Post("/list", r.ctrl.List)
}
