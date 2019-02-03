package models

import (
	"fmt"
	. "github.com/Fallensouls/Pandora/setting"
	. "github.com/Fallensouls/Pandora/util/jsonutil"
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
		DbConfig.Host,
		DbConfig.Port,
		DbConfig.User,
		DbConfig.Name,
		DbConfig.Password))

	if err != nil {
		log.Panic(err)
	}

	engine.ShowSQL(true)
	engine.SetMapper(core.GonicMapper{})

	engine.DB().SetMaxIdleConns(10)
	engine.DB().SetMaxOpenConns(100)
}
