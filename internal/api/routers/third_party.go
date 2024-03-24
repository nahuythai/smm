package routers

import (
	thirdPartyCtrl "smm/internal/api/controller/thirdparty"
	"smm/internal/api/middleware"

	"github.com/gofiber/fiber/v2"
)

type ThirdParty interface {
	V1()
}

type thirdParty struct {
	ctrl   thirdPartyCtrl.Controller
	router fiber.Router
}

func NewThirdParty(router fiber.Router) ThirdParty {
	return &thirdParty{
		ctrl:   thirdPartyCtrl.New(),
		router: router,
	}
}

func (r *thirdParty) V1() {
	router := r.router.Group("/third-party")
	router.Use(middleware.ThirdPartyAuth)
	router.Post("/request", r.ctrl.Route)
}
