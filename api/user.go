package api

import (
	"github.com/Fallensouls/Pandora/errs"
	"github.com/Fallensouls/Pandora/models"
	"github.com/Fallensouls/Pandora/util/jsonutil"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func Register(c *gin.Context) {
	var user models.User
	if c.BindJSON(&user) == nil {
		if err := user.ValidateUserInfo(); err != nil {
			c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
		} else {
			if err := user.AddUser(); err != nil {
				c.Set("error", true)
				if errHandler(err) {
					c.Set("errMsg", err.Error())
				}
			} else {
				c.Status(http.StatusOK)
			}
		}
	}
}

// Login by email address or cellphone number
func Login(c *gin.Context) {
	var user models.User
	if c.BindJSON(&user) == nil {
		if err := user.CheckPassword(); err != nil {
			c.Set("error", true)
			if errHandler(err) {
				c.Set("errMsg", err.Error())
			}
		} else {
			token, err := jsonutil.GenerateJWT(strconv.FormatInt(user.Id, 10))
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
			c.JSON(http.StatusOK, Response{Message: "OK", Data: token})
		}
	}
}

func UpdateProfile(c *gin.Context) {
	var user models.User
	if c.BindJSON(&user) == nil {
		id := c.GetInt64("id")
		if err := user.UpdateUserProfile(id); err != nil {
			c.Set("error", true)
			if errHandler(err) {
				c.Set("errMsg", err.Error())
			}
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func GetProfile(c *gin.Context) {
	id := c.GetInt64("id")
	var user models.User
	if err := user.GetUser(id); err != nil {
		log.Println(err)
		c.Set("error", true)
		if errHandler(err) {
			c.Set("errMsg", err.Error())
		}
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func errHandler(err error) bool {
	switch err {
	case errs.ErrUserNotFound, errs.ErrUserInactive, errs.ErrUserRestricted, errs.ErrUserBanned,
		errs.ErrEmailUsed, errs.ErrCellphoneUsed, errs.ErrEncodingPassword, errs.ErrWrongPassword:
		return true
	default:
		return false
	}
}
