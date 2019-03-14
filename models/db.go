package models

import (
	"fmt"
	. "github.com/go-pandora/core/conf"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"log"
)

type BasicModel struct {
	Id       int64    `json:"id"`
	CreateAt JsonTime `json:"-"    xorm:"created"`
	UpdateAt JsonTime `json:"-"    xorm:"updated"`
}

var engine *xorm.Engine

func init() {
	var err error
	engine, err = xorm.NewEngine("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		Config.DBHost,
		Config.DBPort,
		Config.DBUser,
		Config.DBName,
		Config.DBPassword))

	if err != nil {
		log.Panicln("failed to connect to Postgres:" + err.Error())
	}

	engine.ShowSQL(true)
	engine.SetMapper(core.GonicMapper{})

	engine.DB().SetMaxIdleConns(10)
	engine.DB().SetMaxOpenConns(100)
}
