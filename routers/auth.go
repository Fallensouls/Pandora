package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pandora/core/cache"
	. "github.com/go-pandora/core/conf"
	"github.com/go-pandora/core/errs"
	"github.com/go-pandora/core/models"
	"github.com/go-pandora/pkg/auth/jwt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var auth *jwt.JWTAuth

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func init() {
	var err error
	auth, err = jwt.NewJWTAuth(jwt.Options{
		AccessPublicKey:      []byte(Config.AccessSecret),
		RefreshPublicKey:     []byte(Config.RefreshSecret),
		SigningMethod:        jwt.SetSigningMethod(Config.SigningAlgorithm),
		AccessTokenDuration:  Config.Timeout,
		RefreshTokenDuration: Config.MaxRefreshTime,
	}, cache.IsJWTRevoked, cache.RevokeJWT)

	if err != nil {
		log.Panicln(err)
	}
}

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

	id := strconv.FormatInt(user.Id, 10)
	accessToken, err := auth.CreateAccessToken(id)
	if err != nil {
		err = errs.New(err)
		return
	}
	refreshToken, err := auth.CreateRefreshToken(id)
	if err != nil {
		err = errs.New(err)
		return
	}
	c.JSON(http.StatusOK, Response{Data: gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}})
}

func LogoutByJWT(c *gin.Context) {
	id := c.GetString("user_id")
	if err := auth.Revoke(id); err != nil {
		c.Set("error", errs.New(err))
	}
	c.Status(http.StatusOK)
}

func RefreshToken(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if parts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := parts[1]
	id, err := auth.RefreshChecker(token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	accessToken, err := auth.CreateAccessToken(id.(string))
	if err != nil {
		c.Set("error", errs.New(err))
		return
	}
	c.JSON(http.StatusOK, Response{Data: gin.H{
		"access_token": accessToken,
	},
	})

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
