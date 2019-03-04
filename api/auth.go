package api

import (
	"github.com/Fallensouls/Pandora/cache"
	"github.com/Fallensouls/Pandora/errs"
	"github.com/Fallensouls/Pandora/middleware/jwt"
	"github.com/Fallensouls/Pandora/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

/*
	*************************************
    *     JWT-Based Authentication      *
    *************************************
*/

// Login by email address or cellphone number.
func LoginByJWT(c *gin.Context) {
	var (
		user models.User
		err  error
	)
	defer func() { c.Set("error", err) }()

	if c.BindJSON(&user) != nil {
		return
	}

	if err = user.Login(); err != nil {
		return
	}

	//// A user that has logged in can't login again
	//var status bool
	//status, err = cache.CheckLoginStatus(user.Id)
	//if err != nil{
	//	return
	//}
	//if status{
	//	err = errs.ErrUserLogin
	//	return
	//}

	accessToken, err1 := jwt.GenerateAccessJWT(user.Id)
	refreshToken, err2 := jwt.GenerateRefreshJWT(user.Id)
	err3 := cache.SetJWTDeadline(user.Id)
	if err1 != nil || err2 != nil || err3 != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, Response{Data: gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}})
}

func LogoutByJWT(c *gin.Context) {
	id := c.GetInt64("user_id")
	if err := cache.SetJWTDeadline(id); err != nil {
		c.Set("error", err)
		return
	}

	c.Status(http.StatusOK)
}

func RefreshToken(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
			Message: errs.ErrUnauthenticated.Error(),
		})
		return
	}
	parts := strings.SplitN(auth, " ", 2)
	if parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
			Message: errs.ErrInvalidAuthHeader.Error(),
		})
		return
	}

	token := parts[1]
	if id, timestamp, err := jwt.ValidateRefreshJWT(token); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
			Message: err.Error(),
		})
	} else {
		status, e := cache.CheckJWTInBlacklist(id, timestamp)
		if e != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !status {
			c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
				Message: errs.ErrUserLogout.Error(),
			})
			return
		}
		token, err := jwt.GenerateAccessJWT(id)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, Response{Data: gin.H{
			"access_token": token,
		},
		})
	}
}

/*
	******************************************
	*       Session-Based Authentication     *
    ******************************************
*/

func LoginBySession() {

}

func LogoutBySession() {

}
