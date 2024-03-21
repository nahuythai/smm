package main

import (
	"smm/internal/api/routers"
	"smm/internal/database"
	"smm/pkg/configure"
	"smm/pkg/jwt"
	"smm/pkg/logging"
	"smm/pkg/mail"
	"smm/pkg/otp"
	"smm/pkg/response"
	"smm/pkg/validator"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	cfg    = configure.GetConfig()
	logger = logging.GetLogger()
)

func main() {
	database.InitDatabase()
	validator.InitValidator()
	jwt.New(cfg.SecretKey).InitGlobal()
	otp.New().InitGlobal()
	mail.New(mail.Option{
		MailHost: cfg.MailHost,
		MailPort: cfg.MailPort,
		Email:    cfg.MailEmail,
		Password: cfg.MailPassword,
	}).InitGlobal()
	app := fiber.New(fiber.Config{
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		ErrorHandler: response.FiberErrorHandler,
	})

	addMiddlewares(app)
	addRoute(app)
	if err := app.Listen(cfg.ServerAddr()); err != nil {
		logger.Fatal().Err(err)
	}
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
	routers.NewService(router).V1()
	routers.NewProvider(router).V1()
	routers.NewOrder(router).V1()
}
