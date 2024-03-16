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
	TokenTypeTransaction
)

const (
	LocalTransactionKey = "transaction"
	LocalUserKey        = "user"
)

const (
	UserRoleUser = iota
	UserRoleAdmin
	UserRoleSuperAdmin
)

const (
	TransactionTypeVerifyLogin = iota
	TransactionTypeVerifyEmail
)
