package api

import (
	"github.com/Fallensouls/Pandora/models"
	"github.com/Fallensouls/Pandora/redis"
	"github.com/Fallensouls/Pandora/util/jsonutil"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	var (
		user models.User
		err  error
	)
	defer func() { c.Set("error", err) }()

	if c.BindJSON(&user) != nil {
		return
	}

	if err = user.ValidateUserInfo(); err != nil {
		return
	}

	if err = user.AddUser(); err != nil {
		return
	}

	c.Status(http.StatusOK)
}

// Login by email address or cellphone number
func Login(c *gin.Context) {
	var (
		user models.User
		err  error
	)
	defer func() { c.Set("error", err) }()

	if c.BindJSON(&user) != nil {
		return
	}

	if err = user.CheckPassword(); err != nil {
		return
	}

	//// A user that has logged in can't login again
	//var status bool
	//status, err = redis.CheckLoginStatus(user.Id)
	//if err != nil{
	//	return
	//}
	//if status{
	//	err = errs.ErrUserLogin
	//	return
	//}

	token, err1 := jsonutil.GenerateJWT(user.Id)
	err2 := redis.SetStatusLogin(user.Id)
	if err1 != nil || err2 != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, Response{Message: "OK", Data: token})
}

func Logout(c *gin.Context) {
	id := c.GetInt64("user_id")
	if err := redis.SetStatusLogout(id); err != nil {
		c.Set("error", err)
		return
	}

	c.Status(http.StatusOK)
}

func ActivateUser(c *gin.Context) {
	id := c.GetInt64("id")
	if err := models.ActivateUser(id); err != nil {
		c.Set("error", err)
		return
	}

	c.Status(http.StatusOK)
}

func RestrictUser(c *gin.Context) {
	id := c.GetInt64("id")
	if err := models.RestrictUser(id); err != nil {
		c.Set("error", err)
	}

	c.Status(http.StatusOK)
}

func BanUser(c *gin.Context) {
	id := c.GetInt64("id")
	if err := models.BanUser(id); err != nil {
		c.Set("error", err)
	}

	c.Status(http.StatusOK)
}

func UpdateProfile(c *gin.Context) {
	var user models.User
	if c.BindJSON(&user) != nil {
		return
	}

	id := c.GetInt64("id")
	if err := user.UpdateUserProfile(id); err != nil {
		c.Set("error", err)
		return
	}

	c.Status(http.StatusOK)
}

func GetProfile(c *gin.Context) {
	id := c.GetInt64("id")
	var user models.User

	if err := user.GetUser(id); err != nil {
		c.Set("error", err)
		return
	}

	c.JSON(http.StatusOK, user)
}
