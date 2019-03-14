package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pandora/core/api"
	"github.com/go-pandora/core/errs"
	"net/http"
	"strconv"
)

// SimpleAuthorizer provides a simple authorization for users.
func SimpleAuthorizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			userId := c.GetString("user_id")
			userIdInt64, _ := strconv.ParseInt(userId, 10, 64)
			id := c.GetInt64("id")
			if userIdInt64 != id {
				c.AbortWithStatusJSON(http.StatusForbidden, api.Response{
					Message: errs.ErrUnauthorized.Error(),
				})
				return
			}
		}
	}
}
