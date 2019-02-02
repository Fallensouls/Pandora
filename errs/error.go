package errs

import "errors"

// Define possible errors may happen
var (
	ErrInvalidParam = errors.New("invalid param")
	ErrInvalidData  = errors.New("invalid data")

	ErrInfoRequired     = errors.New("please provide a valid email address or a cellphone number")
	ErrUserNotFound     = errors.New("this account does not exist")
	ErrUserInactive     = errors.New("please activate your account first")
	ErrUserRestricted   = errors.New("this account has been restricted")
	ErrUserBanned       = errors.New("this account has been banned")
	ErrWrongPassword    = errors.New("this password is wrong")
	ErrEncodingPassword = errors.New("fail to encode your password")
	ErrEmailUsed        = errors.New("this email address has already been used")
	ErrCellphoneUsed    = errors.New("this cellphone number has already been used")
)
