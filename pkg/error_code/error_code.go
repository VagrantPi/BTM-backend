package error_code

const (
	// user
	ErrInvalidJWTParse = "ErrInvalidJWTParse"
	ErrInvalidJWT      = "ErrInvalidJWT"
	ErrJWT             = "ErrJWT"

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
)
