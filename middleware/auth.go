package middleware

import (
	"github.com/Fallensouls/Pandora/api"
	"github.com/Fallensouls/Pandora/errs"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
