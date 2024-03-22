package constants

const (
	ErrMsgResourceNotFound = "resource not found"
	ErrMsgFieldWrongType   = "field wrong type"
	ErrMsgOtpWrong         = "wrong otp"
)

const (
	ErrCodeAppUnknown             = "0000"
	ErrCodeAppInternalServerError = "0001"
	ErrCodeAppBadRequest          = "0002"
	ErrCodeAppForbidden           = "0003"
	ErrCodeAppUnauthorized        = "0004"
)

const (
	ErrCodeCategoryNotFound = "0100"
)

const (
	ErrCodeUserNotFound                = "0200"
	ErrCodeUserExist                   = "0201"
	ErrCodeUserBanned                  = "0202"
	ErrCodeUserInvalidEmail            = "0203"
	ErrCodeUserAdminPermissionRequired = "0204"
	ErrCodeUserNotEnoughBalance        = "0205"
)

const (
	ErrCodeSessionNotFound = "0300"
)

const (
	ErrCodeOtpNotFound = "0400"
	ErrCodeOtpWrong    = "0401"
)

const (
	ErrCodeTokenWrong = "0500"
)

const (
	ErrCodeServiceNotFound = "0600"
)

const (
	ErrCodeProviderNotFound = "0700"
)

const (
	ErrCodeOrderNotFound        = "0800"
	ErrCodeOrderQuantityInvalid = "0801"
)

const (
	ErrBackgroundTaskExist = "0900"
)

const (
	ErrCodeCustomRateNotFound = "1000"
)

const (
	ErrCodeTransactionNotFound = "1100"
)

const (
	ErrCodePaymentMethodNotFound = "1200"
)
