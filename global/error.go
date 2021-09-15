package global

const (
	SuccessCode      = 200
	UnauthorizedCode = 401
	NotFoundCode     = 404
	SystemErrorCode  = 500
)

const (
	UnknownCacheErrorCode = 5001
	UnknownDBErrorCode    = 5002
)

const (
	IllegalInputErrCode = 4001
)

const (
	AccountPasswordErrCode = 4031
)

const (
	NotFoundMsg           = "404 not found"
	SuccessMsg            = "success"
	UnauthorizedMsg       = "unauthorized, please re-authorize"
	ServerErrorMsg        = "server internal error"
	IllegalInputErrMsg    = "input is illegal"
	AccountPasswordErrMsg = "incorrect account or password"
	EmptyTokenErrMsg      = "token not obtained"
	TokenHasExpiredMsg    = "token has expired"
)

type Error interface {
	//Code error code returned to the client
	Code() int
	//Msg error message returned to the client
	Msg() string
}

type ResponseError struct {
	ErrCode int
	Reason  string
}

func (re ResponseError) Code() int {
	return re.ErrCode
}

func (re ResponseError) Msg() string {
	return re.Reason
}
