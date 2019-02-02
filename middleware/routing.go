package middleware

import (
	"github.com/Fallensouls/Pandora/api"
	"github.com/Fallensouls/Pandora/errs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// IdValidator validates whether id is valid
// with GET, PUT, DELETE method, we expect that id is integer(int64)
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

// Errhandler checks whether there is an error
// if not, it will do nothing and just return
func Errhandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if errExist := c.GetBool("error"); errExist == true {
			if errMsg := c.GetString("errMsg"); errMsg != "" {
				c.JSON(http.StatusBadRequest, api.Response{Message: errMsg})
			} else {
				c.Status(http.StatusInternalServerError)
			}
		}
		return
	}
}
