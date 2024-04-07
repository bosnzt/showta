package msg

import (
	"errors"
)

var (
	ErrMustLogin    = errors.New("errMustLogin")
	ErrTokenExpired = errors.New("errTokenExpired")
	ErrTokenInvalid = errors.New("errTokenInvalid")
	ErrAuthAccount  = errors.New("errAuthAccount")
	ErrAccessPwd    = errors.New("errAccessPwd")
)
