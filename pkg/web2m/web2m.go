package web2m

import (
	"smm/pkg/logging"

	"github.com/gofiber/fiber/v2"
)

var (
	logger = logging.GetLogger()
)

type Service interface {
	GetTransactionHistory() (*TransactionHistoryResponse, error)
	GetTransactionHistoryWithUrl(url string) (*TransactionHistoryResponse, error)
}
type service struct {
	url string
}

func New(url string) Service {
	return &service{
		url: url,
	}
}

const (
	TransactionTypeIn  = "IN"
	TransactionTypeOut = "OUT"
)

type TransactionHistoryResponse struct {
	Status       bool          `json:"status"`
	Message      string        `json:"message"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	TransactionId   string `json:"transactionID"`
	Amount          string `json:"amount"`
	Description     string `json:"description"`
	TransactionDate string `json:"transactionDate"`
	Type            string `json:"TYPE"`
}

func (s *service) GetTransactionHistory() (*TransactionHistoryResponse, error) {
	var res TransactionHistoryResponse
	code, body, errs := fiber.AcquireClient().Get(s.url).Struct(&res)
	logger.Debug().Caller().Str("func", "GetTransactionHistory").Str("funcInline", "fiber.AcquireClient").Int("code", code).Bytes("body", body).Msg("web2m-api")
	if len(errs) > 0 {
		logger.Error().Errs("errs", errs).Caller().Str("func", "GetTransactionHistory").Str("funcInline", "fiber.AcquireClient").Msg("web2m-api")
		return nil, errs[0]
	}
	return &res, nil
}

func (s *service) GetTransactionHistoryWithUrl(url string) (*TransactionHistoryResponse, error) {
	var res TransactionHistoryResponse
	code, body, errs := fiber.AcquireClient().Get(url).Struct(&res)
	logger.Debug().Caller().Str("func", "GetTransactionHistory").Str("funcInline", "fiber.AcquireClient").Int("code", code).Bytes("body", body).Msg("web2m-api")
	if len(errs) > 0 {
		logger.Error().Errs("errs", errs).Caller().Str("func", "GetTransactionHistoryWithUrl").Str("funcInline", "fiber.AcquireClient").Msg("web2m-api")
		return nil, errs[0]
	}
	return &res, nil
}
