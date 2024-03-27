package routers

import (
	transactionCtrl "smm/internal/api/controller/transaction"
	"smm/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

type Transaction interface {
	V1()
}

type transaction struct {
	ctrl   transactionCtrl.Controller
	router fiber.Router
}

func NewTransaction(router fiber.Router) Transaction {
	return &transaction{
		ctrl:   transactionCtrl.New(),
		router: router,
	}
}

func (r *transaction) V1() {
	r.Admin()
	r.User()
}

func (r *transaction) Admin() {
	router := r.router.Group("/admin/transactions")
	router.Use(middleware.UserAuth, middleware.AdminPermission)
	router.Post("/list", r.ctrl.List)
}

func (r *transaction) User() {
	router := r.router.Group("/transactions")
	router.Use(middleware.UserAuth)
	router.Post("/list", r.ctrl.UserList)
}
