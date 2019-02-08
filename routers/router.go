// Package routers specifies routing for each API.
package routers

import (
	"github.com/Fallensouls/Pandora/api"
	"github.com/Fallensouls/Pandora/middleware"
	. "github.com/Fallensouls/Pandora/setting"
	"github.com/gin-gonic/gin"
)

func SetRouter() (r *gin.Engine) {
	r = gin.Default()
	gin.SetMode(Config.RunMode)
	r.Use(middleware.ErrHandler())

	r.POST("/register", api.Register)
	r.POST("/login", api.Login)
	r.PUT("/activate/:id", api.ActivateUser)
	r.PUT("/logout", middleware.Authenticator(), api.Logout)

	Api := r.Group("/api")
	Api.Use(middleware.IdValidator(), middleware.Authenticator(), middleware.SimpleAuthorizer())
	{
		Api.GET("/user/:id", api.GetProfile)
		Api.PUT("/user/:id", api.UpdateProfile)
	}

	return
}
