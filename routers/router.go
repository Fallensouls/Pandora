package routers

import (
	"Pandora/controller"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"log"
	"net/http"
	"time"
)

type ServerConfig struct {
	RunMode			string			`ini:"run_mode"`
	Port			string			`ini:"http_port"`
	ReadTimeout		time.Duration	`ini:"read_timeout"`
	WriteTimeout	time.Duration 	`ini:"write_timeout"`
}

func setServerConfig(c *ServerConfig)  {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil{
		log.Panic("fail to load config file")
	}
	if err = cfg.Section("server").MapTo(c); err != nil{
		log.Panic("fail to set server config")
	}
	c.RunMode = cfg.Section("").Key("run_mode").String()
}

func SetRouter() (r *gin.Engine, config ServerConfig)  {
	setServerConfig(&config)
	r = gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(config.RunMode)

	test := r.Group("/test")
	{
		test.GET("", func(c *gin.Context) {
			c.String(http.StatusOK, "hello world")
		})
	}

	user := r.Group("/user")
	{
		user.POST("/register", controller.Register)
		user.POST("/login", controller.Login)
		user.GET("/:id", controller.GetProfile)
		user.PUT("/:id", controller.UpdateProfile)
	}

	return
}

