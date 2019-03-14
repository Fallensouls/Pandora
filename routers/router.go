// Package routers specifies routing for each API.
package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pandora/core/api"
	. "github.com/go-pandora/core/conf"
	"github.com/go-pandora/core/middleware"
)

func SetRouter() (r *gin.Engine) {
	r = gin.Default()
	gin.SetMode(Config.RunMode)
	r.Use(middleware.ErrHandler())

	r.MaxMultipartMemory = 4 << 20
	Upload := r.Group("/upload")
	Upload.Use(auth.Authenticator())
	{
		Upload.POST("/avatar", api.UploadAvatar)
	}

	Auth := r.Group("/auth")
	{
		Auth.POST("/register", api.Register)
		Auth.POST("/login", LoginByJWT)
		//Auth.GET("/activate", api.ActivateUser)
		Auth.PUT("/logout", auth.Authenticator(), LogoutByJWT)
		Auth.GET("/refresh", RefreshToken)
	}

	Api := r.Group("/api")
	Api.Use(middleware.IdValidator(), auth.Authenticator(), middleware.SimpleAuthorizer())
	{
		Api.GET("/user/:id", api.GetProfile)
		Api.PUT("/user/:id", api.UpdateProfile)
	}

	return
}
