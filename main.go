package main

import (
	"log"
	"smm/internal/api/routers"
	"smm/internal/database"
	"smm/pkg/configure"
	"smm/pkg/logging"
	"smm/pkg/response"
	"smm/pkg/validator"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var cfg = configure.GetConfig()

func main() {
	database.InitDatabase()
	validator.InitValidator()
	app := fiber.New(fiber.Config{
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		ErrorHandler: response.FiberErrorHandler,
	})

	addMiddlewares(app)
	addRoute(app)
	log.Fatal(app.Listen(cfg.ServerAddr()))
}

func addMiddlewares(app *fiber.App) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(logging.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
}

func addRoute(app fiber.Router) {
	router := app.Group("/api/v1")
	routers.NewCategory(router).V1()
	routers.NewUser(router).V1()
}
