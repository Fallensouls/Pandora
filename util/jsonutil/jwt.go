package jsonutil

import (
	. "github.com/Fallensouls/Pandora/setting"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTClaims struct {
	jwt.StandardClaims
	Identity string `json:"identity"`
}

// GenerateJWT generates JsonWebToken used for authentication.
// Identity is whatever you consider as the unique id for user.
func GenerateJWT(identity string) (token string, err error) {
	claim := JWTClaims{
		Identity: identity,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(JwtConfig.Duration).Unix(),
			Issuer:    "Fallensouls",
			NotBefore: time.Now().Unix(),
		},
	}
	unsigned := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	token, err = unsigned.SignedString(JwtConfig.Secret)
	return
}

//func ValidateJWT(token string, lastModify time.Time) (result bool, err error) {
//
//}
