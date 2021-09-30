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
	IllegalInputErrCode      = 4001
	RegisteredErrCode        = 4002
	NoUserOrPasswordCode     = 4003
	PluginNameExistCode      = 4004
	InvalidUserCode          = 4005
	NoRecodeCode             = 4006
	PluginVersionIsExistCode = 4007
)

const (
	AccountPasswordErrCode = 4031
)

const (
	NotFoundMsg             = "404 not found"
	SuccessMsg              = "success"
	UnauthorizedMsg         = "unauthorized, please re-authorize"
	ServerErrorMsg          = "server internal error"
	IllegalInputErrMsg      = "input is illegal"
	AccountPasswordErrMsg   = "incorrect account or password"
	EmptyTokenErrMsg        = "token not obtained"
	TokenHasExpiredMsg      = "token has expired"
	RegisteredErrMsg        = "input info is registered"
	NoUserOrPasswordMsg     = "no user or password mismatch "
	PluginNameExistMsg      = "already exist the same name plugin"
	InvalidUserMsg          = "invalid user"
	NoRecodeMsg             = "no corresponding data"
	PluginVersionIsExistMsg = "plugin version number is existed"
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

func NewResponseSystemError() ResponseError {
	return ResponseError{ErrCode: SystemErrorCode, Reason: ServerErrorMsg}
}
