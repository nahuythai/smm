package constants

const (
	CategoryStatusOn = iota
	CategoryStatusOff
)

const (
	UserStatusActive = iota
	UserStatusBanned
)

const (
	UserLanguageEnglish = iota
	UserLanguageVietnamese
)

const (
	TokenTypeAccess = iota
	TokenTypeSession
)

const (
	LocalSessionKey = "session"
	LocalUserKey    = "user"
)

const (
	UserRoleUser = iota
	UserRoleAdmin
	UserRoleSuperAdmin
)

const (
	SessionTypeVerifyLogin = iota
	SessionTypeVerifyEmail
	SessionTypeQRTopUpPayment
)

const (
	ServiceStatusOn = iota
	ServiceStatusOff
)

const (
	ProviderStatusOn = iota
	ProviderStatusOff
)

const (
	BackgroundTypeUpdateBalance = iota
)

const (
	TransactionTypePlayOrder = iota
	TransactionTypeAddBalance
)

const (
	PaymentMethodStatusOn = iota
	PaymentMethodStatusOff
)

const (
	PaymentStatusPending = iota
	PaymentStatusCompleted
	PaymentStatusCancelled
)

const (
	ThirdPartyActionOrderCreate         = "add"
	ThirdPartyActionOrderStatus         = "status"
	ThirdPartyActionOrderMultipleStatus = "orders"
	ThirdPartyActionServiceList         = "services"
	ThirdPartyActionUserBalance         = "balance"
)

const (
	OrderStatusAwaiting = iota
	OrderStatusPending
	OrderStatusProcessing
	OrderStatusInProgress
	OrderStatusCompleted
	OrderStatusPartial
	OrderStatusCanceled
	OrderStatusRefunded
)

var OrderStatusTextMapping = map[int]string{
	OrderStatusAwaiting:   "AWAITING",
	OrderStatusPending:    "PENDING",
	OrderStatusProcessing: "PROCESSING",
	OrderStatusInProgress: "INPROGRESS",
	OrderStatusCompleted:  "COMPLETED",
	OrderStatusPartial:    "PARTIAL",
	OrderStatusCanceled:   "CANCELED",
	OrderStatusRefunded:   "REFUNDED",
}

var OrderStatusMapping = map[string]int{
	"AWAITING":   OrderStatusAwaiting,
	"PENDING":    OrderStatusPending,
	"PROCESSING": OrderStatusProcessing,
	"INPROGRESS": OrderStatusInProgress,
	"COMPLETED":  OrderStatusCompleted,
	"PARTIAL":    OrderStatusPartial,
	"CANCELED":   OrderStatusCanceled,
	"REFUNDED":   OrderStatusRefunded,
}
