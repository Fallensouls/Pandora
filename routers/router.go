// Package routers specifies routing for each API.
package routers

import (
	"github.com/Fallensouls/Pandora/api"
	. "github.com/Fallensouls/Pandora/conf"
	"github.com/Fallensouls/Pandora/middleware"
	"github.com/gin-gonic/gin"
)

func SetRouter() (r *gin.Engine) {
	r = gin.Default()
	gin.SetMode(Config.RunMode)
	r.Use(middleware.ErrHandler())

	r.MaxMultipartMemory = 4 << 20
	upload := r.Group("/upload")
	upload.Use(middleware.Authenticator())
	{
		upload.POST("/avatar", api.UploadAvatar)
	}

	auth := r.Group("/auth")
	{
		auth.POST("/register", api.Register)
		auth.POST("/login", api.Login)
		auth.GET("/activate", api.ActivateUser)
		auth.PUT("/logout", middleware.Authenticator(), api.Logout)
		auth.GET("/refresh", api.RefreshToken)
	}

	Api := r.Group("/api")
	Api.Use(middleware.IdValidator(), middleware.Authenticator(), middleware.SimpleAuthorizer())
	{
		Api.GET("/user/:id", api.GetProfile)
		Api.PUT("/user/:id", api.UpdateProfile)
	}

	return
}
