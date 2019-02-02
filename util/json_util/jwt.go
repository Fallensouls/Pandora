package json_util

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-ini/ini"
	"log"
	"time"
)

type JWTConfig struct {
	Secret   string        `ini:"secret"`
	Duration time.Duration `ini:"duration"`
}

type JWTClaims struct {
	jwt.StandardClaims
	Identity string `json:"identity"`
}

var jwtConfig = &JWTConfig{}

func init() {
	cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Panic("fail to load config file")
	}
	if err = cfg.Section("jwt").MapTo(jwtConfig); err != nil {
		log.Println("fail to set jwt config, use default jwt config...")
		jwtConfig.Duration = time.Hour
		jwtConfig.Secret = "Hatsune Miku"
	}
}

func GenerateJWT(identity string) (token string, err error) {
	claim := JWTClaims{
		Identity: identity,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtConfig.Duration).Unix(),
			Issuer:    "Fallensouls",
			NotBefore: time.Now().Unix(),
		},
	}
	unsigned := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	token, err = unsigned.SignedString(jwtConfig.Secret)
	return
}

//func ValidateJWT(token string, lastModify time.Time) (result bool, err error) {
//
//}
