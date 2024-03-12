package routers

import (
	categoryCtrl "smm/internal/api/controller/category"

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
	router := r.router.Group("/categories")
	router.Post("/", r.ctrl.Create)
	router.Post("/list", r.ctrl.List)
	router.Put("/", r.ctrl.Update)
	router.Get("/:id", r.ctrl.Get)
}
