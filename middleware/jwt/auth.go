package jwt

import (
	"github.com/Fallensouls/Pandora/api"
	"github.com/Fallensouls/Pandora/cache"
	"github.com/Fallensouls/Pandora/errs"
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
		if id, timestamp, err := ValidateAccessJWT(token); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.Response{
				Message: err.Error(),
			})
		} else {
			status, e := cache.CheckJWTInBlacklist(id, timestamp)
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
