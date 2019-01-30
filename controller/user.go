package controller

import (
	"Pandora/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Register(c *gin.Context)  {
	var user models.User
	if c.BindJSON(&user) == nil{
		if by, err := user.ValidateUserInfo(); err != nil{
			c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
		}else {
			if err := models.AddUser(user, by); err != nil {
				c.JSON(http.StatusInternalServerError, Response{Message: err.Error()})
			}else {
				c.Status(http.StatusOK)
			}
		}
	}
}

// Login by email address or cellphone number
func Login(c *gin.Context)  {
	var user models.User
	if c.BindJSON(&user) == nil{
		if err := models.ComparePassword(user); err != nil{
			if err == models.ErrWrongPassword || err == models.ErrUserAbnormal{
				c.JSON(http.StatusUnauthorized, Response{Message: err.Error()})
			} else {
				c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
			}
		} else{
			c.Status(http.StatusOK)
		}
	}
}

func UpdateProfile(c *gin.Context)  {
	var user models.User
	if c.BindJSON(&user) == nil{
		if err := models.UpdateUserProfile(user); err != nil{
			c.Status(http.StatusInternalServerError)
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func GetProfile(c *gin.Context)  {
	if id, err := strconv.ParseInt(c.Param("id"), 10 ,64); err != nil{
		c.Status(http.StatusBadRequest)
	} else{
		if user, err := models.GetUser(id); err != nil{
			if err == models.ErrUserNotFound{
				c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
			}
			if err == models.ErrUserAbnormal{
				c.JSON(http.StatusUnauthorized, Response{Message: err.Error()})
			}
			c.Status(http.StatusInternalServerError)
		}else {
			c.JSON(http.StatusOK, user)
		}
	}
}
