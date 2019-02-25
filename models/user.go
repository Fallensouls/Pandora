package models

import (
	"github.com/Fallensouls/Pandora/errs"
	. "github.com/Fallensouls/Pandora/util/jsonutil"
	. "github.com/Fallensouls/Pandora/util/valiutil"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	BasicModel  `xorm:"extends"`
	Username    string      `json:"username"`
	Password    string      `json:"password,omitempty"`
	Avatar      []byte      `json:"avatar,omitempty" xorm:"-"`
	Age         int         `json:"age,omitempty"`
	Gender      int         `json:"gender,omitempty"`
	Address     string      `json:"address,omitempty"`
	Description string      `json:"description,omitempty"`
	Email       *string     `json:"email,omitempty"`
	Cellphone   *string     `json:"cellphone,omitempty"`
	Auth        []Authority `json:"-" xorm:"extends"`
	Status      int         `json:"-"`
	LastModify  JsonTime    `json:"-"`
}

type Authority struct {
	BasicModel `xorm:"extends"`
	Role       string
	User       []User `xorm:"extends"`
}

// Define user's status
const (
	Inactive   = 0 // sign up but has not been activated
	Normal     = 1 //
	Restricted = 2 // login is permitted, but can only receive read-only message
	Banned     = 3 // can't login
)

// Gender
const (
	Unknown = 0
	Male    = 1
	Female  = 2
)

// TableName specifies the table name of struct User
func (u *User) TableName() string {
	return "users"
}

// GetUser provides user's information.
// Note that password, email address and cellphone number won't return.
func (u *User) GetUser(id int64) error {
	if exist, err := engine.ID(id).Omit("password", "email", "cellphone").Get(u); err != nil {
		return errs.New(err)
	} else {
		if !exist {
			return errs.ErrUserNotFound
		}
	}
	switch u.Status {
	case Inactive:
		return errs.ErrUserInactive
	case Restricted:
		return errs.ErrUserRestricted
	case Banned:
		return errs.ErrUserBanned
	default:
		return nil
	}
}

// AddUser will add a new user.
// Users can sign up by email address or cellphone number.
func (u *User) AddUser() error {
	if err := encodePassword(&u.Password); err != nil {
		return errs.New(err)
	}
	if _, err := engine.Insert(u); err != nil {
		if strings.Contains(err.Error(), "email") {
			return errs.ErrEmailUsed
		} else if strings.Contains(err.Error(), "cellphone") {
			return errs.ErrCellphoneUsed
		} else {
			return errs.New(err)
		}
	}
	return nil
}

// UpdateUserProfile will update user's profile.
func (u *User) UpdateUserProfile(id int64) error {
	if _, err := engine.ID(id).Cols("age", "gender", "address", "description").
		Update(u); err != nil {
		return errs.New(err)
	}
	return nil
}

// CheckPassword compares password provided by user and password stored in database.
func (u *User) CheckPassword() error {
	pw := u.Password
	u.Password = ""
	if exist, err := engine.Cols("id", "password", "status").Get(u); err != nil {
		return errs.New(err)
	} else {
		if !exist {
			return errs.ErrUserNotFound
		}
	}
	switch u.Status {
	case Inactive:
		return errs.ErrUserInactive
	case Banned:
		return errs.ErrUserBanned
	default:
		break
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw)); err != nil {
		return errs.ErrWrongPassword
	}
	return nil
}

func (u *User) ChangeEmail(id int64) error {
	if _, err := engine.ID(id).Cols("email").Update(u); err != nil {
		return errs.New(err)
	}
	return nil
}

func (u *User) ChangeCellphone(id int64) error {
	if _, err := engine.ID(id).Cols("cellphone").Update(u); err != nil {
		return errs.New(err)
	}
	return nil
}

func ActivateUser(id int64) error {
	return changeStatus(id, Normal)
}

func RestrictUser(id int64) error {
	return changeStatus(id, Restricted)
}

func BanUser(id int64) error {
	return changeStatus(id, Banned)
}

// ValidateUserInfo validates whether user's information is valid.
func (u *User) ValidateUserInfo() (err error) {
	if u.Email == nil && u.Cellphone == nil {
		return errs.ErrInfoRequired
	}
	if err = ValidateUsername(u.Username); err != nil {
		return errs.ErrInvalidUsername
	}
	if err = ValidatePassword(u.Password); err != nil {
		return errs.ErrInvalidPassword
	}
	if u.Email != nil {
		if err = ValidateEmail(*u.Email); err != nil {
			return errs.ErrInvalidEmail
		}
	}
	if u.Cellphone != nil {
		if err = ValidateCellphone(*u.Cellphone); err != nil {
			return errs.ErrInvalidCellphone
		}
	}
	return
}

// EncodePassword will encode password with bcrypt.
func encodePassword(password *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return errs.ErrEncodingPassword
	}
	*password = string(hash)
	return nil
}

func ChangePassword(id int64, old string, new string) error {
	var user User
	if err := user.GetUser(id); err != nil {
		return errs.New(err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(old), []byte(user.Password)); err != nil {
		return errs.ErrWrongPassword
	}
	if err := encodePassword(&user.Password); err != nil {
		return errs.ErrEncodingPassword
	}
	var modifyTime JsonTime
	modifyTime.GetJsonTime()
	if _, err := engine.ID(user.Id).Cols("last_modify").
		Update(&User{LastModify: modifyTime}); err != nil {
		return errs.New(err)
	}
	return nil
}

func changeStatus(id int64, status int) error {
	if _, err := engine.ID(id).Cols("status").Update(&User{Status: status}); err != nil {
		return errs.New(err)
	}
	return nil
}
