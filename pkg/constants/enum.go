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
	OrderStatusAwaiting:   "Awaiting",
	OrderStatusPending:    "Pending",
	OrderStatusProcessing: "Processing",
	OrderStatusInProgress: "InProgress",
	OrderStatusCompleted:  "Completed",
	OrderStatusPartial:    "Partial",
	OrderStatusCanceled:   "Canceled",
	OrderStatusRefunded:   "Refunded",
}
