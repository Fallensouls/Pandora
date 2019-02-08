package middleware

import (
	"github.com/Fallensouls/Pandora/api"
	"github.com/Fallensouls/Pandora/errs"
	"github.com/Fallensouls/Pandora/redis"
	"github.com/Fallensouls/Pandora/util/jsonutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Authenticator checks whether user is authenticated.
func Authenticator() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			return
		}
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.Response{
				Message: errs.ErrUnauthenticated.Error(),
			})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.Response{
				Message: errs.ErrInvalidAuthHeader.Error(),
			})
			return
		}

		token := parts[1]
		if id, timestamp, err := jsonutil.ValidateJWT(token); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.Response{
				Message: err.Error(),
			})
		} else { // check if a user really logs in
			status, e := redis.CheckJWTStatus(id, timestamp)
			if e != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			if !status {
				c.AbortWithStatusJSON(http.StatusUnauthorized, api.Response{
					Message: errs.ErrUserLogout.Error(),
				})
				return
			}
			c.Set("user_id", id)
		}
	}
}

// SimpleAuthorizer provides a simple authorization for users.
func SimpleAuthorizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			userId := c.GetInt64("user_id")
			id := c.GetInt64("id")
			if userId != id {
				c.AbortWithStatusJSON(http.StatusBadRequest, api.Response{
					Message: errs.ErrUnauthorized.Error(),
				})
				return
			}
		}
	}
}
