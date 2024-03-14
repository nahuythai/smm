package routers

import (
	userCtrl "smm/internal/api/controller/user"
	"smm/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

type User interface {
	V1()
}

type user struct {
	ctrl   userCtrl.Controller
	router fiber.Router
}

func NewUser(router fiber.Router) User {
	return &user{
		ctrl:   userCtrl.New(),
		router: router,
	}
}

func (r *user) V1() {
	router := r.router.Group("/users")
	router.Post("/", r.ctrl.Create)
	router.Post("/list", r.ctrl.List)
	router.Put("/", r.ctrl.Update)
	router.Get("/:id", r.ctrl.Get)
	router.Post("/:id/generate-api-key", r.ctrl.GenerateApiKey)
	router.Post("/update-balance", r.ctrl.UpdateBalance)
	router.Post("/update-password", r.ctrl.UpdatePassword)
	router.Post("/login", r.ctrl.Login)
	transactionRouter := router.Group("/transactions")
	transactionRouter.Use(middleware.TransactionAuth)
	transactionRouter.Post("/login-verify", r.ctrl.VerifyLogin)
}
