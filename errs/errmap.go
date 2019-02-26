package errs

var ErrMap = map[string]error{
	"1001": ErrInvalidParam,
	"1002": ErrInvalidData,
	"1003": ErrInvalidToken,

	"20001": ErrInfoRequired,
	"20002": ErrInvalidUsername,
	"20003": ErrInvalidPassword,
	"20004": ErrInvalidEmail,
	"20005": ErrInvalidCellphone,
	"20006": ErrUserNotFound,
	"20007": ErrUserInactive,
	"20008": ErrUserRestricted,
	"20009": ErrUserBanned,
	"20010": ErrWrongPassword,
	"20011": ErrEncodingPassword,
	"20012": ErrEmailUsed,
	"20013": ErrCellphoneUsed,
	"20014": ErrUserLogin,
	"20015": ErrUserLogout,
}
