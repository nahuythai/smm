package cron

import (
	"context"
	"fmt"
	"smm/internal/database/models"
	"smm/internal/database/queries"
	"smm/pkg/configure"
	"smm/pkg/constants"
	"smm/pkg/logging"
	"smm/pkg/providerapi"
	"smm/pkg/request"
	"smm/pkg/web2m"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var runner *cron.Cron

var (
	cfg    = configure.GetConfig()
	logger = logging.GetLogger()
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
	// runner.AddFunc("*/1 * * * *", updateOrdersStatus)
	// updateOrdersStatus()
	// updatePaymentStatus()
}

func updateOrdersStatus() {
	logger.Info().Str("func", "updateOrdersStatus").Str("info", "task start!").Msg("cron")
	defer func() {
		logger.Info().Str("func", "updateOrdersStatus").Str("info", "task done!").Msg("cron")
	}()
	var limit int64 = 1000
	var page int64 = 0

	//get all key
	providerQueryOption := queries.NewOption()
	providerQueryOption.SetOnlyField("api_key", "_id", "url")
	providers, err := queries.NewProvider(backgroundCtx).GetByFilter(bson.M{"status": constants.ProviderStatusOn}, providerQueryOption)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "updateOrdersStatus").
			Str("funcInline", "queries.NewProvider(backgroundCtx).GetByFilter").Msg("cron")
		return
	}
	providerIdMap := make(map[primitive.ObjectID]models.Provider)
	for _, provider := range providers {
		providerIdMap[provider.Id] = provider
	}

	orderFilter := bson.M{
		"status": bson.M{
			"$in": []int{
				constants.OrderStatusAwaiting, constants.OrderStatusPending, constants.OrderStatusProcessing,
				constants.OrderStatusInProgress, constants.OrderStatusPartial},
		},
	}

	for {
		page += 1
		queryOption := queries.NewOption()
		pagination := request.NewPagination(page, limit)
		queryOption.SetPagination(pagination)
		queryOption.SetOnlyField("provider_order_id", "created_at", "_id", "provider_id")
		orderQuery := queries.NewOrder(backgroundCtx)
		orders, err := orderQuery.GetByFilter(orderFilter, queryOption)
		if err != nil {
			logger.Error().Err(err).Caller().Str("func", "updateOrdersStatus").
				Str("funcInline", "orderQuery.GetByFilter").Msg("cron")
		}
		if len(orders) < int(limit) && page > 1 {
			return
		}

		orderWrites := make([]mongo.WriteModel, 0, len(orders))
		providerOrderIdMap := make(map[primitive.ObjectID][]int64)
		// TODO: update transaction and balance by result
		// transactionWrites := make([]mongo.WriteModel, 0, len(orders))
		for _, order := range orders {
			if order.ProviderOrderId == 0 && time.Since(order.CreatedAt) > 5*time.Minute {
				orderWrites = append(orderWrites, mongo.NewUpdateOneModel().
					SetFilter(bson.M{"_id": order.Id}).SetUpdate(bson.M{"$set": bson.M{"status": constants.OrderStatusRefunded}}))
			} else {
				providerOrderIdMap[order.ProviderId] = append(providerOrderIdMap[order.ProviderId], order.ProviderOrderId)
			}
		}

		// call to provider
		for providerId, provider := range providerIdMap {
			orderIds, ok := providerOrderIdMap[providerId]
			if !ok {
				continue
			}
			orderStrIds := lo.Map(orderIds, func(orderId int64, index int) string {
				return fmt.Sprintf("%d", orderId)
			})
			providerApi := providerapi.New(provider.Url, provider.ApiKey)
			res, err := providerApi.MultipleOrderStatus(providerapi.MultipleOrderStatusRequest{
				Orders: strings.Join(orderStrIds, ","),
			})
			if err != nil {
				logger.Error().Err(err).Caller().Str("func", "updateOrdersStatus").
					Str("funcInline", "providerApi.MultipleOrderStatus").Msg("cron")
				continue
			}
			for _, result := range res {
				if result.Error == "" {
					orderWrites = append(orderWrites, mongo.NewUpdateOneModel().
						SetFilter(bson.M{"provider_order_id": result.Order}).
						SetUpdate(bson.M{"$set": bson.M{
							"status":        constants.OrderStatusMapping[strings.ToUpper(result.Status)],
							"remains":       result.Remains,
							"start_counter": result.StartCounter,
							"updated_at":    time.Now(),
						}}))
				} else {
					orderWrites = append(orderWrites, mongo.NewUpdateOneModel().
						SetFilter(bson.M{"provider_order_id": result.Order}).
						SetUpdate(bson.M{"$set": bson.M{
							"updated_at":              time.Now(),
							"provider_order_response": result.Error,
						}}))
				}
			}
		}
		if err = orderQuery.BulkWrite(orderWrites); err != nil {
			logger.Error().Err(err).Caller().Str("func", "updateOrdersStatus").
				Str("funcInline", "orderQuery.BulkWrite").Msg("cron")
		}
	}
}

func updatePaymentStatus() {
	logger.Info().Str("func", "updatePaymentStatus").Str("info", "task start!").Msg("cron")
	defer func() {
		logger.Info().Str("func", "updatePaymentStatus").Str("info", "task done!").Msg("cron")
	}()
	var limit int64 = 1000
	var page int64 = 0
	listPaymentHistoryApi := []string{
		"https://api.web2m.com/historyapimbv3/PasswordMBBank/AccountNumber/TokenMBBank",
	}
	transactionIdMap := make(map[string]web2m.Transaction)
	for _, url := range listPaymentHistoryApi {
		res, err := web2m.New(url).GetTransactionHistory()
		if err != nil {
			logger.Error().Err(err).Caller().Str("func", "updatePaymentStatus").
				Str("funcInline", "web2m.New(url).GetTransactionHistory").Msg("cron")
			continue
		}
		if res.Status {
			for _, transaction := range res.Transactions {
				if transaction.Type == web2m.TransactionTypeIn {
					transactionIdMap[transaction.Description] = transaction
				}
			}
		}
	}

	for {
		page += 1
		queryOption := queries.NewOption()
		pagination := request.NewPagination(page, limit)
		queryOption.SetPagination(pagination)
		queryOption.SetOnlyField("amount", "_id", "created_at", "trx_id")
		paymentQuery := queries.NewPayment(backgroundCtx)
		payments, err := paymentQuery.GetByFilter(bson.M{"status": constants.PaymentStatusPending}, queryOption)
		if err != nil {
			logger.Error().Err(err).Caller().Str("func", "updatePaymentStatus").
				Str("funcInline", "paymentQuery.GetByFilter").Msg("cron")
		}
		if len(payments) < int(limit) && page > 1 {
			return
		}
		paymentWrites := make([]mongo.WriteModel, 0, len(payments))
		for _, payment := range payments {
			if transaction, ok := transactionIdMap[payment.TransactionId]; ok {
				transactionAmount, _ := strconv.ParseFloat(transaction.Amount, 64)
				if payment.Amount == transactionAmount {
					paymentWrites = append(paymentWrites, mongo.NewUpdateOneModel().
						SetFilter(bson.M{"_id": payment.Id}).
						SetUpdate(bson.M{"$set": bson.M{
							"status":     constants.PaymentStatusCompleted,
							"updated_at": time.Now(),
						}}))
				}
			} else {
				if time.Since(payment.CreatedAt) > 24*3*time.Hour {
					paymentWrites = append(paymentWrites, mongo.NewUpdateOneModel().
						SetFilter(bson.M{"_id": payment.Id}).
						SetUpdate(bson.M{"$set": bson.M{
							"status":     constants.PaymentStatusCancelled,
							"updated_at": time.Now(),
							"feedback":   "cron expired",
						}}))
				}
			}
		}
		if err = paymentQuery.BulkWrite(paymentWrites); err != nil {
			logger.Error().Err(err).Caller().Str("func", "updateOrdersStatus").
				Str("funcInline", "orderQuery.BulkWrite").Msg("cron")
		}
	}
}
