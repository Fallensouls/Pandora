// Package routers specifies routing for each API.
package routers

import (
	"github.com/Fallensouls/Pandora/api"
	"github.com/Fallensouls/Pandora/middleware"
	. "github.com/Fallensouls/Pandora/setting"
	"github.com/gin-gonic/gin"
)

func SetRouter() (r *gin.Engine) {
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(ServerCfg.RunMode)

	root := r.Group("")
	root.Use(middleware.Errhandler())
	{
		root.POST("/register", api.Register)
		root.POST("/login", api.Login)
	}

	Api := r.Group("/api")
	Api.Use(middleware.IdValidator())
	Api.Use(middleware.Errhandler())
	{
		Api.GET("/user/:id", api.GetProfile)
		Api.PUT("/user/:id", api.UpdateProfile)
	}

	return
}
