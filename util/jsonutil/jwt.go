package jsonutil

import (
	"fmt"
	"github.com/Fallensouls/Pandora/errs"
	. "github.com/Fallensouls/Pandora/setting"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTClaims struct {
	jwt.StandardClaims
	Id int64 `json:"id"`
}

var jwtSecret = []byte(Config.Secret)

// GenerateJWT generates Json Web Token used for authentication.
// Here we use user's id as extra data.
// Please do not add important information such as password to payload of JWT.
func GenerateJWT(id int64) (token string, err error) {
	claim := JWTClaims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(Config.Timeout).Unix(),
			Issuer:    Config.Issuer,
			NotBefore: time.Now().Unix(),
		},
	}

	unsigned := jwt.NewWithClaims(getSigningMethod(Config.SigningAlgorithm), claim)
	token, err = unsigned.SignedString(jwtSecret)
	return
}

// ValidateJWT validates whether jwt is valid.
// If so, we still have to check if user really logged in before.
func ValidateJWT(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	claims := token.Claims.(jwt.MapClaims)
	if claims["id"] == nil || err != nil {
		return 0, errs.ErrInvalidToken
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return 0, errs.ErrInvalidToken
	}
	return int64(id), nil
}

func getSigningMethod(method string) *jwt.SigningMethodHMAC {
	switch method {
	case "HS256":
		return jwt.SigningMethodHS256
	case "HS384":
		return jwt.SigningMethodHS384
	case "HS512":
		return jwt.SigningMethodHS512
	default:
		return jwt.SigningMethodHS256
	}
}
