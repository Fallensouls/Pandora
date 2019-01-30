package models

import (
	"errors"
	"github.com/Fallensouls/Pandora/util/date"
	"github.com/Fallensouls/Pandora/util/validation"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type User struct {
	BasicModel
	Username    string       `json:"username"`
	Password    string       `json:"password,omitempty"`
	Avatar      *uuid.UUID   `json:"avatar,omitempty"`
	Age         int          `json:"age,omitempty"`
	Gender      int          `json:"gender,omitempty"`
	Address     string       `json:"address,omitempty"`
	Description string       `json:"description,omitempty"`
	Email       string       `json:"email,omitempty"`
	Cellphone   string       `json:"cellphone,omitempty"`
	Auth        *[]Authority `json:"-"`
	Status      int          `json:"-"`
	LastLogin   time.Time    `json:"-"         gorm:"column:lastlogin"`
	LastModify  time.Time    `json:"-"         gorm:"column:lastmodify"`
}

type Authority struct {
	BasicModel
	Role string
	User *[]User
}

// Define possible errors may happen
var (
	ErrUserNotFound     = errors.New("the user does not exist")
	ErrUserAbnormal     = errors.New("the user is under abnormal mode")
	ErrWrongPassword    = errors.New("the password is wrong")
	ErrEncodingPassword = errors.New("fail to encode the password")
	ErrEmailUsed        = errors.New("the email address has already been used")
	ErrCellphoneUsed    = errors.New("the cellphone number has already been used")
	ErrChangeStatus     = errors.New("fail to change status")
)

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

// Password, email address and cellphone number won't return to frontend
func GetUser(id int64) (user User, err error) {
	if err = db.Find(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrUserNotFound
			return
		}
	}
	if user.Status == Inactive || user.Status == Banned {
		err = ErrUserAbnormal
		return
	}
	user.Password = ""
	user.Email = ""
	user.Cellphone = ""
	return
}

// by represents sign up by email address or cellphone number
func AddUser(user User, by string) (err error) {
	err = SetPassword(&user.Password)
	if err != nil {
		return ErrEncodingPassword
	}
	if err = db.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), by) {
			switch by {
			case "email":
				return ErrEmailUsed
			case "cellphone":
				return ErrCellphoneUsed
			default:
				return
			}
		}
	}
	return
}

// Update user's profile
func UpdateUserProfile(user User) error {
	if err := db.Model(&user).Select("avatar,age,gender,address,description").
		Updates(user).Error; err != nil {
		return err
	}
	return nil
}

// Encode password with bcrypt
func SetPassword(password *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*password = string(hash)
	return nil
}

func ChangePassword(id int64, old string, new string) error {
	user, err := GetUser(id)
	if err != nil {
		return err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(old), []byte(user.Password)); err != nil {
		return ErrWrongPassword
	}
	err = SetPassword(&user.Password)
	if err != nil {
		return ErrEncodingPassword
	}
	user.LastModify = date.GetStandardTime()
	if err = db.Model(&user).Select("password, lastmodify").Updates(user).Error; err != nil {
		return err
	}
	return nil
}

// compare user's password and password stored in database
func ComparePassword(user User) error {
	pw := user.Password
	if err := db.Select("password, status").Find(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUserNotFound
		}
		if user.Status == Inactive || user.Status == Banned {
			return ErrUserAbnormal
		}
	}
	if err := bcrypt.CompareHashAndPassword([]byte(pw), []byte(user.Password)); err != nil {
		return ErrWrongPassword
	}
	return nil
}

func ChangeEmail(user User) error {
	//user := User{BasicModel: BasicModel{ID:id}, Email: email}
	if err := db.Model(&user).Select("email").Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func ChangeCellphone(user User) error {
	//user := User{BasicModel: BasicModel{ID:id}, Cellphone: cell}
	if err := db.Model(&user).Select("cellphone").Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func RestrictUser(id int64) error {
	return ChangeStatus(id, Restricted)
}

func BanUser(id int64) error {
	return ChangeStatus(id, Banned)
}

func ChangeStatus(id int64, status int) error {
	user := User{BasicModel: BasicModel{ID: id}, Status: status}
	if err := db.Model(&user).Select("status").Updates(user).Error; err != nil {
		return ErrChangeStatus
	}
	return nil
}

func (u *User) ValidateUserInfo() (by string, err error) {
	if u.Email == "" && u.Cellphone == "" {
		err = errors.New("please provide a valid email address or a cellphone number")
		return
	}
	if err = validation.ValidateUsername(u.Username); err != nil {
		return
	}
	if err = validation.ValidatePassword(u.Password); err != nil {
		return
	}
	if u.Email != "" {
		by = "email"
		if err = validation.ValidateEmail(u.Email); err != nil {
			return
		}
	}
	if u.Cellphone != "" {
		by = "cellphone"
		if err = validation.ValidateCellphone(u.Cellphone); err != nil {
			return
		}
	}
	return
}
