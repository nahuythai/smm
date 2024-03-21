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
	r.User()
	r.Admin()
}

func (r *user) User() {
	userRouter := r.router.Group("/users")
	userRouter.Post("/login", r.ctrl.Login)
	userRouter.Post("/register", r.ctrl.Register)
	userRouter.Get("/verify-email", r.ctrl.VerifyEmail)

	sessionRouter := userRouter.Group("/sessions")
	sessionRouter.Use(middleware.SessionAuth)
	sessionRouter.Post("/login-verify", r.ctrl.VerifyLogin)

	userRouter.Use(middleware.UserAuth)
	userRouter.Get("/me", r.ctrl.Me)
}

func (r *user) Admin() {
	adminRouter := r.router.Group("/admin/users")
	adminRouter.Use(middleware.UserAuth, middleware.AdminPermission)
	adminRouter.Post("/", r.ctrl.Create)
	adminRouter.Post("/list", r.ctrl.List)
	adminRouter.Put("/", r.ctrl.Update)
	adminRouter.Get("/:id", r.ctrl.Get)
	adminRouter.Post("/:id/generate-api-key", r.ctrl.GenerateApiKey)
	adminRouter.Post("/update-balance", r.ctrl.UpdateBalance)
	adminRouter.Post("/update-password", r.ctrl.UpdatePassword)
}
