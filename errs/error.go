// Package errs defines all possible errors may occur in runtime.
package errs

type Err struct {
	SystemError bool
	Message     string
}

func (e *Err) Error() string {
	return e.Message
}

func New(err error) error {
	return &Err{true, err.Error()}
}

var (
	ErrInvalidParam = &Err{Message: "invalid param"}
	ErrInvalidData  = &Err{Message: "invalid data"}
	ErrInvalidToken = &Err{Message: "invalid token"}
)

var (
	ErrInvalidAuthHeader = &Err{Message: "your auth header is invalid"}
	ErrUnauthenticated   = &Err{Message: "please login"}
	ErrUnauthorized      = &Err{Message: "you are not authorized"}
)

var (
	ErrInfoRequired     = &Err{Message: "please provide a valid email address or a cellphone number"}
	ErrInvalidUsername  = &Err{Message: "your username is not valid"}
	ErrInvalidPassword  = &Err{Message: "your password is not valid"}
	ErrInvalidEmail     = &Err{Message: "your email address is not valid"}
	ErrInvalidCellphone = &Err{Message: "your cellphone number is not valid"}
	ErrUserNotFound     = &Err{Message: "this account does not exist"}
	ErrUserInactive     = &Err{Message: "please activate your account first"}
	ErrUserRestricted   = &Err{Message: "this account has been restricted"}
	ErrUserBanned       = &Err{Message: "this account has been banned"}
	ErrWrongPassword    = &Err{Message: "incorrect password"}
	ErrEncodingPassword = &Err{Message: "failed to encode your password"}
	ErrEmailUsed        = &Err{Message: "this email address has already been used"}
	ErrCellphoneUsed    = &Err{Message: "this cellphone number has already been used"}
	ErrUserLogin        = &Err{Message: "you have logged in"}
	ErrUserLogout       = &Err{Message: "you have logged out"}
)
