package error_code

const (
	// user
	ErrInvalidJWTParse = "ErrInvalidJWTParse"
	ErrInvalidJWT      = "ErrInvalidJWT"
	ErrJWT             = "ErrJWT"
	ErrTokenExpired    = "ErrTokenExpired"

	// whitelist
	ErrWhitelistDuplicate = "ErrWhitelistDuplicate"

	// system error
	ErrInvalidRequest     = "ErrInvalidRequest"
	ErrDBError            = "ErrDBError"
	ErrRedisError         = "ErrRedisError"
	ErrDiError            = "ErrDiError"
	ErrInternalError      = "ErrInternalError"
	ErrInternalPanicError = "ErrInternalPanicError"
	ErrForbidden          = "ErrForbidden"
	ErrThirdPartyHttpCall = "ErrThirdPartyHttpCall"

	// cib
	ErrCibTokenParse = "ErrCibTokenParse"
	ErrCibTokenFetch = "ErrCibTokenFetch"

	// third party
	ErrRemoveFile = "ErrRemoveFile"
	ErrCreateFile = "ErrCreateFile"

	// tools
	ErrUnzipFile = "ErrUnzipFile"
	ErrCsvOpen   = "ErrCsvOpen"
	ErrCsvRead   = "ErrCsvRead"
)
