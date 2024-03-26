package providerapi

import (
	"smm/pkg/logging"

	"github.com/gofiber/fiber/v2"
)

var (
	logger = logging.GetLogger()
)

type Provider interface {
	AddOrder(request AddOrderRequest) (*AddOrderResponse, error)
	MultipleOrderStatus(request MultipleOrderStatusRequest) ([]MultipleOrderStatusOrderResponse, error)
}
type provider struct {
	apiKey      string
	providerUrl string
}

func New(providerUrl, apiKey string) Provider {
	return &provider{
		apiKey:      apiKey,
		providerUrl: providerUrl,
	}
}

type AddOrderResponse struct {
	OrderId int64  `json:"order"`
	Status  string `json:"status"`
	Error   string `json:"error"`
}

type AddOrderRequest struct {
	ApiServiceId string `json:"api_service_id"`
	Link         string `json:"link"`
	Quantity     int64  `json:"quantity"`
}

type MultipleOrderStatusRequest struct {
	Orders string `json:"orders"`
}

type MultipleOrderStatusOrderResponse struct {
	Status       string  `json:"status"`
	StartCounter int64   `json:"start_counter"`
	Remains      int64   `json:"remains"`
	Charge       *string `json:"charge"`
	Order        int64   `json:"order"`
	Error        string  `json:"error"`
}

func (p *provider) AddOrder(request AddOrderRequest) (*AddOrderResponse, error) {
	var res AddOrderResponse
	code, body, errs := fiber.AcquireClient().Post(p.providerUrl).JSON(fiber.Map{
		"key":      p.apiKey,
		"action":   "add",
		"service":  request.ApiServiceId,
		"link":     request.Link,
		"quantity": request.Quantity,
	}).Struct(&res)
	logger.Debug().Caller().Str("func", "AddOrder").Str("funcInline", "fiber.AcquireClient").Int("code", code).Bytes("body", body).Msg("provider-api")
	if len(errs) > 0 {
		logger.Error().Errs("errs", errs).Caller().Str("func", "AddOrder").Str("funcInline", "fiber.AcquireClient").Msg("provider-api")
		return nil, errs[0]
	}
	return &res, nil
}

func (p *provider) MultipleOrderStatus(request MultipleOrderStatusRequest) ([]MultipleOrderStatusOrderResponse, error) {
	var res []MultipleOrderStatusOrderResponse
	code, body, errs := fiber.AcquireClient().Post(p.providerUrl).JSON(fiber.Map{
		"key":    p.apiKey,
		"action": "orders",
		"orders": request.Orders,
	}).Struct(&res)
	logger.Debug().Caller().Str("func", "MultipleOrderStatus").Str("funcInline", "fiber.AcquireClient").Int("code", code).Bytes("body", body).Msg("provider-api")
	if len(errs) > 0 {
		logger.Error().Errs("errs", errs).Caller().Str("func", "MultipleOrderStatus").Str("funcInline", "fiber.AcquireClient").Msg("provider-api")
		return nil, errs[0]
	}
	return res, nil
}
