// Package conf only accomplishes loading configuration.
package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Configuration struct {
	*Database
	*Server
	*Redis
	*JWT
	*Storage
}

type Database struct {
	//Type     string
	DBName     string `yaml:"name"`
	DBUser     string `yaml:"user"`
	DBPassword string `yaml:"password"`
	DBHost     string `yaml:"host"`
	DBPort     string `yaml:"port"`
}

type Server struct {
	RunMode      string
	Port         string
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type Redis struct {
	RedisHost     string `yaml:"host"`
	RedisPort     string `yaml:"port"`
	RedisPassword string `yaml:"password"`
}

type JWT struct {
	SigningAlgorithm string        `yaml:"signing_algorithm"`
	AccessSecret     string        `yaml:"access_secret"`
	RefreshSecret    string        `yaml:"refresh_secret"`
	Timeout          time.Duration `yaml:"duration"`
	Issuer           string
	MaxRefreshTime   time.Duration `yaml:"max_refresh_time"`
}

type Storage struct {
	*Image
}

type Image struct {
	AvatarPath string `yaml:"avatar"`
}

var Config Configuration

func init() {
	loadConfig()
	checkDatabase()
	checkServer()
	checkRedis()
	checkJWT()
	checkStorage()
}

func loadConfig() {
	data, err := ioutil.ReadFile("conf/config.yaml")
	if err != nil {
		log.Panicln("failed to load config file")
	}
	if err = yaml.Unmarshal(data, &Config); err != nil {
		log.Panicln("failed to load configuration")
	}
}

func checkDatabase() {
	if Config.Database == nil {
		log.Panicln("failed to init database configuration")
	}
}

func checkServer() {
	if Config.Server == nil {
		log.Panicln("failed to init server configuration")
	}
}

func checkRedis() {
	if Config.Redis == nil {
		log.Panicln("failed to init cache configuration")
	}
}

func checkJWT() {
	if Config.JWT == nil {
		log.Println("failed to init jwt configuration, use default jwt config...")
		Config.JWT = &JWT{}
		Config.Timeout = time.Hour
		Config.MaxRefreshTime = 7 * 24 * time.Hour
		Config.AccessSecret = "Hatsune Miku"
		Config.RefreshSecret = "Miku-chan maji tenshi"
		Config.Issuer = "Fallensouls"
	} else {
		Config.Timeout *= time.Minute
		Config.MaxRefreshTime *= time.Hour
	}
}

func checkStorage() {
	if Config.Storage == nil {
		log.Panicln("failed to init storage configuration")
	}

	if err := os.MkdirAll(Config.AvatarPath, os.ModeDir); err != nil {
		log.Panicf("failed to create avatar path, error: %s", err.Error())
	}

}
