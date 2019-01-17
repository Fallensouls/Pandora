package models

import (
	"github.com/go-ini/ini"
	"time"
)

type config struct {
	Type         string				`ini:"type"`
	Name         string				`ini:"name"`
	User         string				`ini:"user"`
	Password     string				`ini:"password"`
	Host         string				`ini:"host"`
}

type BaseModel struct {
	ID				int64 			`json:"id"`
	CreateOn		time.Time		`json:"-"`
	UpdateOn		time.Time		`json:"-"`
}

//var db *gorm.DB
func init()  {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil{

	}

	c := new(config)
	err = cfg.Section("database").MapTo(c)

}