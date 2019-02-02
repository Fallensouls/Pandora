package models

import (
	"fmt"
	. "github.com/Fallensouls/Pandora/util/json_util"
	"github.com/go-ini/ini"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"log"
)

type DatabaseConfig struct {
	//Type     string `ini:"type"`
	Name     string `ini:"name"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	Host     string `ini:"host"`
	Port     string `ini:"port"`
}

type BasicModel struct {
	Id       int64    `json:"id"`
	CreateAt JsonTime `json:"-"    xorm:"created"`
	UpdateAt JsonTime `json:"-"    xorm:"updated"`
}

var engine *xorm.Engine

func init() {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Panic("fail to load config file")
	}
	c := new(DatabaseConfig)
	if err = cfg.Section("database").MapTo(c); err != nil {
		log.Panic("fail to set database config")
	}
	engine, err = xorm.NewEngine("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		c.Host,
		c.Port,
		c.User,
		c.Name,
		c.Password))

	if err != nil {
		log.Panic(err)
	}

	engine.ShowSQL(true)
	engine.SetMapper(core.GonicMapper{})

	engine.DB().SetMaxIdleConns(10)
	engine.DB().SetMaxOpenConns(100)
}
