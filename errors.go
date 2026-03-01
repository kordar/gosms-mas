package mas

import "github.com/kordar/gosms"

var masErrorMap = map[string]gosms.ErrorCode{
	"success": gosms.ErrSuccess,

	"InvalidUsrOrPwd": gosms.ErrAuthFailed,
	"IllegalMac":      gosms.ErrAuthFailed,

	"IllegalSignId": gosms.ErrSignInvalid,
	"NoSignId":      gosms.ErrSignInvalid,

	"InvalidMessage": gosms.ErrInvalidRequest,

	"TooManyMobiles": gosms.ErrTooManyMobiles,
}

func mapMASError(rspcod string) gosms.ErrorCode {
	if c, ok := masErrorMap[rspcod]; ok {
		return c
	}
	return gosms.ErrUnknown
}
