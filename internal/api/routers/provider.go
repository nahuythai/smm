package routers

import (
	providerCtrl "smm/internal/api/controller/provider"
	"smm/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

type Provider interface {
	V1()
}

type provider struct {
	ctrl   providerCtrl.Controller
	router fiber.Router
}

func NewProvider(router fiber.Router) Provider {
	return &provider{
		ctrl:   providerCtrl.New(),
		router: router,
	}
}

func (r *provider) V1() {
	r.Admin()
}

func (r *provider) Admin() {
	router := r.router.Group("/admin/providers")
	router.Use(middleware.UserAuth, middleware.AdminPermission)
	router.Post("/", r.ctrl.Create)
	router.Post("/list", r.ctrl.List)
	router.Put("/", r.ctrl.Update)
	router.Get("/:id", r.ctrl.Get)
	router.Delete("/:id", r.ctrl.Delete)
}
