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
	OrderStatusAwaiting = iota
	OrderStatusPending
	OrderStatusProcessing
	OrderStatusInProgress
	OrderStatusCompleted
	OrderStatusPartial
	OrderStatusCanceled
	OrderStatusRefunded
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
