package error_code

const (
	// btm user
	ErrInvalidJWTParse = "ErrInvalidJWTParse"
	ErrInvalidJWT      = "ErrInvalidJWT"
	ErrJWT             = "ErrJWT"
	ErrTokenExpired    = "ErrTokenExpired"

	// btm whitelist
	ErrWhitelistDuplicate = "ErrWhitelistDuplicate"

	// btm sumsub
	ErrBTMSumsubGetItem          = "ErrBTMSumsubGetItem"
	ErrBTMSumsubCreateItem       = "ErrBTMSumsubCreateItem"
	ErrBTMSumsubIdNumberNotFound = "ErrBTMSumsubIdNumberNotFound"

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

	// sumsub
	ErrSumsubRequest      = "ErrSumsubRequest"
	ErrSumsubBadRequest   = "ErrSumsubBadRequest"
	ErrSumsubApiUnmarshal = "ErrSumsubApiUnmarshal"

	// tools
	ErrToolsUnzipFile         = "ErrToolsUnzipFile"
	ErrToolsCsvOpen           = "ErrToolsCsvOpen"
	ErrToolsCsvRead           = "ErrToolsCsvRead"
	ErrToolsHttpRequest       = "ErrToolsHttpRequest"
	ErrToolsHttpRequestDo     = "ErrToolsHttpRequestDo"
	ErrToolsHttpRequestIoRead = "ErrToolsHttpRequestIoRead"
	ErrToolsHttpMarshal       = "ErrToolsHttpMarshal"
)
