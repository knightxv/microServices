package constant

import "errors"

// key = errCode, string = errMsg
type ErrInfo struct {
	ErrCode int32
	ErrMsg  string
}

var (
	OK        = ErrInfo{0, ""}
	ErrServer = ErrInfo{500, "server error"}

	ErrParseToken = ErrInfo{700, ErrParseTokenMsg.Error()}

	ErrTokenExpired     = ErrInfo{701, ErrTokenExpiredMsg.Error()}
	ErrTokenInvalid     = ErrInfo{702, ErrTokenInvalidMsg.Error()}
	ErrTokenNotValidYet = ErrInfo{704, ErrTokenNotValidYetMsg.Error()}
	ErrTokenMalformed   = ErrInfo{703, ErrTokenMalformedMsg.Error()}
	ErrTokenUnknown     = ErrInfo{705, ErrTokenUnknownMsg.Error()}

	ErrAccess   = ErrInfo{ErrCode: 801, ErrMsg: ErrAccessMsg.Error()}
	ErrStatus   = ErrInfo{ErrCode: 804, ErrMsg: ErrStatusMsg.Error()}
	ErrDB       = ErrInfo{ErrCode: 802, ErrMsg: ErrDBMsg.Error()}
	ErrArgs     = ErrInfo{ErrCode: 803, ErrMsg: ErrArgsMsg.Error()}
	ErrCallback = ErrInfo{ErrCode: 809, ErrMsg: ErrCallBackMsg.Error()}
)

var (
	ErrParseTokenMsg       = errors.New("parse token failed")
	ErrTokenExpiredMsg     = errors.New("token is timed out, please log in again")
	ErrTokenInvalidMsg     = errors.New("token has been invalidated")
	ErrTokenNotValidYetMsg = errors.New("token not active yet")
	ErrTokenMalformedMsg   = errors.New("that's not even a token")
	ErrTokenUnknownMsg     = errors.New("couldn't handle this token")
	ErrAccessMsg           = errors.New("no permission")
	ErrStatusMsg           = errors.New("status is abnormal")
	ErrDBMsg               = errors.New("db failed")
	ErrArgsMsg             = errors.New("args failed")
	ErrCallBackMsg         = errors.New("callback failed")

	ErrThirdPartyMsg = errors.New("third party error")
)

const (
	NoError              = 0
	FormattingError      = 10001
	HasRegistered        = 10002
	NotRegistered        = 10003
	PasswordErr          = 10004
	GetTokenErr          = 10005
	RepeatSendCode       = 10006
	MailSendCodeErr      = 10007
	SmsSendCodeErr       = 10008
	CodeInvalidOrExpired = 10009
	RegisterFailed       = 10010
	ResetPasswordFailed  = 10011
	DatabaseError        = 10002
	ServerError          = 10004
	HttpError            = 10005
	IoError              = 10006
	IntentionalError     = 10007
)

func (e ErrInfo) Error() string {
	return e.ErrMsg
}

func (e *ErrInfo) Code() int32 {
	return e.ErrCode
}
