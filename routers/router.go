package routers

import (
	"github.com/Fallensouls/Pandora/api"
	"github.com/Fallensouls/Pandora/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"log"
	"time"
)

type ServerConfig struct {
	RunMode      string        `ini:"run_mode"`
	Port         string        `ini:"http_port"`
	ReadTimeout  time.Duration `ini:"read_timeout"`
	WriteTimeout time.Duration `ini:"write_timeout"`
}

var Server = &ServerConfig{}

func init() {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Panic("fail to load config file")
	}
	if err = cfg.Section("server").MapTo(Server); err != nil {
		log.Println(err)
		log.Panic("fail to set server config")
	}
	Server.RunMode = cfg.Section("").Key("run_mode").String()
}

func SetRouter() (r *gin.Engine) {
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(Server.RunMode)

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
