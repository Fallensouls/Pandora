// Package middleware includes all gin-style middleware.
// Both API and middleware are in handlers chain.
package middleware

import (
	"github.com/Fallensouls/Pandora/api"
	"github.com/Fallensouls/Pandora/errs"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// IdValidator validates whether id is valid.
// With regard to GET, PUT and DELETE, we expect that id is integer(int64).
func IdValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		var (
			id  int64
			err error
		)
		if method == "POST" {
			return
		}
		if id, err = strconv.ParseInt(c.Param("id"), 10, 64); err != nil {
			switch method {
			case "GET", "PUT", "DELETE":
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": errs.ErrInvalidParam.Error(),
				})
			default:
				break
			}
		} else {
			c.Set("id", id)
		}
	}
}

// ErrHandler checks whether there is an error after API is called.
// If not, it will do nothing and just return
func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if err, ok := c.Get("error"); ok {
			if err == nil {
				return
			}
			e := err.(*errs.Err)
			if e.SystemError {
				log.Println(e.Message)
				c.Status(http.StatusInternalServerError)
			} else {
				c.JSON(http.StatusBadRequest, api.Response{Message: e.Message})
			}
		}
	}
}
