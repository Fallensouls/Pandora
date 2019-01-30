package models

import (
	"Pandora/util/date"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"time"
)

type DatabaseConfig struct {
	Type         string				`ini:"type"`
	Name         string				`ini:"name"`
	User         string				`ini:"user"`
	Password     string				`ini:"password"`
	Host         string				`ini:"host"`
	Port         string				`ini:"port"`
}

type BasicModel struct {
	ID				int64 			`json:"id"`
	CreateOn		time.Time		`json:"-"    gorm:"column:createon"`
	UpdateOn		time.Time		`json:"-"    gorm:"column:updateon"`
}

var db *gorm.DB

func init() {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil{
		log.Panic("fail to load config file")
	}
	c := new(DatabaseConfig)
	if err = cfg.Section("database").MapTo(c); err != nil{
		log.Panic("fail to set database config")
	}
	db, err = gorm.Open(c.Type, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		c.Host,
		c.Port,
		c.User,
		c.Name,
		c.Password))

	if err != nil {
		log.Panic(err)
	}

	//gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
	//	return tablePrefix + defaultTableName;
	//}

	//db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.Callback().Create().Replace("gorm:update_time_stamp", createTimeCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeCallback)
}

func CloseDB() {
	defer db.Close()
}


func createTimeCallback(scope *gorm.Scope){
	if !scope.HasError(){
		now := date.GetStandardTime()
		scope.SetColumn("CreateOn",now)
		scope.SetColumn("UpdateOn",now)
	}
}

func updateTimeCallback(scope *gorm.Scope){
	if !scope.HasError(){
		scope.SetColumn("UpdateOn", date.GetStandardTime())
	}
}