// Package setting only accomplishes loading configuration.
package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

type DatabaseConfig struct {
	//Type     string `ini:"type"`
	Name     string `ini:"name"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Host     string `ini:"host"`
	Port     string `ini:"port"`
}

type ServerConfig struct {
	RunMode      string        `ini:"run_mode"`
	Port         string        `ini:"http_port"`
	ReadTimeout  time.Duration `ini:"read_timeout"`
	WriteTimeout time.Duration `ini:"write_timeout"`
}

type JWTConfig struct {
	Secret   string        `ini:"secret"`
	Duration time.Duration `ini:"duration"`
}

var (
	DbConfig  = &DatabaseConfig{}
	ServerCfg = &ServerConfig{}
	JwtConfig = &JWTConfig{}
)

func init() {
	if cfg, err := ini.Load("conf/app.ini"); err != nil {
		log.Panic("fail to load config file")
	} else {
		if err := cfg.Section("database").MapTo(DbConfig); err != nil {
			log.Panic("fail to set database config")
		}

		if err := cfg.Section("server").MapTo(ServerCfg); err != nil {
			log.Panic("fail to set server config")
		}
		ServerCfg.RunMode = cfg.Section("").Key("run_mode").String()

		if err := cfg.Section("jwt").MapTo(JwtConfig); err != nil {
			log.Println("fail to set jwt config, use default jwt config...")
			JwtConfig.Duration = time.Hour
			JwtConfig.Secret = "Hatsune Miku"
		}
	}
}
