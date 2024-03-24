package cron

import (
	"context"
	"fmt"
	"smm/pkg/configure"

	"github.com/robfig/cron/v3"
)

var runner *cron.Cron

var (
	cfg = configure.GetConfig()
)

var (
	backgroundCtx, _ = context.WithTimeout(context.Background(), cfg.MongoRequestTimeout)
)

func Run() {
	runner = cron.New()
	route()
	runner.Start()
}

func route() {
	runner.AddFunc("*/1 * * * *", updateOrdersStatus)
}

func updateOrdersStatus() {
	// providerapi.New(cfg)
	// orders, err := queries.NewOrder(backgroundCtx).GetByFilter()
	fmt.Println("hihi test")
}
