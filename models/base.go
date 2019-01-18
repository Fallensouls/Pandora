package models

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type DatabaseConfig struct {
	Type         string				`ini:"type"`
	Name         string				`ini:"name"`
	User         string				`ini:"user"`
	Password     string				`ini:"password"`
	Host         string				`ini:"host"`
}

type BaseModel struct {
	ID				int64 			`json:"id"`
	CreateOn		time.Time		`json:"-"    time_format:"2006-01-02  15:04:05"`
	UpdateOn		time.Time		`json:"-"    time_format:"2006-01-02  15:04:05"`
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
	db, err = gorm.Open(c.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Name))

	if err != nil {
		log.Println(err)
	}

	//gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
	//	return tablePrefix + defaultTableName;
	//}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CloseDB() {
	defer db.Close()
}