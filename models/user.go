package models

import (
	. "github.com/Fallensouls/Pandora/errs"
	. "github.com/Fallensouls/Pandora/util/json_util"
	. "github.com/Fallensouls/Pandora/util/validate"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	BasicModel  `xorm:"extends"`
	Username    string      `json:"username"`
	Password    string      `json:"password,omitempty"`
	Avatar      *uuid.UUID  `json:"avatar,omitempty"`
	Age         int         `json:"age,omitempty"`
	Gender      int         `json:"gender,omitempty"`
	Address     string      `json:"address,omitempty"`
	Description string      `json:"description,omitempty"`
	Email       *string     `json:"email,omitempty"`
	Cellphone   *string     `json:"cellphone,omitempty"`
	Auth        []Authority `json:"-" xorm:"extends"`
	Status      int         `json:"-"`
	LastLogin   JsonTime    `json:"-"`
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

// GetUser will return user's information.
// Note that password, email address and cellphone number won't return
func (u *User) GetUser(id int64) (err error) {
	var exist bool
	if exist, err = engine.ID(id).Omit("password", "email", "cellphone").Get(u); err != nil {
		return
	} else {
		if !exist {
			err = ErrUserNotFound
			return
		}
	}
	switch u.Status {
	case Inactive:
		err = ErrUserInactive
	case Restricted:
		err = ErrUserRestricted
	case Banned:
		err = ErrUserBanned
	default:
		break
	}
	return
}

// AddUser will add a new user.
// user can sign up by email address or cellphone number
func (u *User) AddUser() (err error) {
	err = encodePassword(&u.Password)
	if err != nil {
		return
	}
	if _, err = engine.Insert(u); err != nil {
		if strings.Contains(err.Error(), "email") {
			return ErrEmailUsed
		} else if strings.Contains(err.Error(), "cellphone") {
			return ErrCellphoneUsed
		} else {
			return
		}
	}
	return
}

// UpdateUserProfile will update user's profile
func (u *User) UpdateUserProfile(id int64) error {
	if _, err := engine.ID(id).Cols("avatar", "age", "gender", "address", "description").
		Update(u); err != nil {
		return err
	}
	return nil
}

// CheckPassword compare user's password and password stored in database
func (u *User) CheckPassword() error {
	pw := u.Password
	u.Password = ""
	if exist, err := engine.Cols("id", "password", "status").Get(u); err != nil {
		return err
	} else {
		if !exist {
			return ErrUserNotFound
		}
	}
	switch u.Status {
	case Inactive:
		return ErrUserInactive
	case Banned:
		return ErrUserBanned
	default:
		break
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw)); err != nil {
		return ErrWrongPassword
	}
	var loginTime JsonTime
	loginTime.GetJsonTime()
	if _, err := engine.ID(u.Id).Cols("last_login").
		Update(&User{LastLogin: loginTime}); err != nil {
		return err
	}
	return nil
}

func (u *User) ChangeEmail(id int64) error {
	if _, err := engine.ID(id).Cols("email").Update(u); err != nil {
		return err
	}
	return nil
}

func (u *User) ChangeCellphone(id int64) error {
	if _, err := engine.ID(id).Cols("cellphone").Update(u); err != nil {
		return err
	}
	return nil
}

// ValidateUserInfo validates whether user's information is valid
func (u *User) ValidateUserInfo() (err error) {
	if u.Email == nil && u.Cellphone == nil {
		return ErrInfoRequired
	}
	if err = ValidateUsername(u.Username); err != nil {
		return ErrInvalidData
	}
	if err = ValidatePassword(u.Password); err != nil {
		return ErrInvalidData
	}
	if u.Email != nil {
		if err = ValidateEmail(*u.Email); err != nil {
			return ErrInvalidData
		}
	}
	if u.Cellphone != nil {
		if err = ValidateCellphone(*u.Cellphone); err != nil {
			return ErrInvalidData
		}
	}
	return
}

// encodePassword will encode password with bcrypt
func encodePassword(password *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return ErrEncodingPassword
	}
	*password = string(hash)
	return nil
}

func ChangePassword(id int64, old string, new string) error {
	var user User
	if err := user.GetUser(id); err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(old), []byte(user.Password)); err != nil {
		return ErrWrongPassword
	}
	if err := encodePassword(&user.Password); err != nil {
		return ErrEncodingPassword
	}
	var modifyTime JsonTime
	modifyTime.GetJsonTime()
	if _, err := engine.ID(user.Id).Cols("last_modify").
		Update(&User{LastModify: modifyTime}); err != nil {
		return err
	}
	return nil
}

func RestrictUser(id int64) error {
	return changeStatus(id, Restricted)
}

func BanUser(id int64) error {
	return changeStatus(id, Banned)
}

func changeStatus(id int64, status int) error {
	if _, err := engine.ID(id).Cols("status").Update(&User{Status: status}); err != nil {
		return err
	}
	return nil
}
