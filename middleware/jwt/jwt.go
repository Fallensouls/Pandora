package jwt

import (
	"fmt"
	. "github.com/Fallensouls/Pandora/conf"
	"github.com/Fallensouls/Pandora/errs"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTClaims struct {
	jwt.StandardClaims
	Id int64 `json:"id"`
}

var (
	accessSecret  = []byte(Config.AccessSecret)
	refreshSecret = []byte(Config.RefreshSecret)
)

// generateJWT generates Json Web Token used for authentication.
// Here we use user's id as extra data.
// Please do not add important information such as password to payload of JWT.
func generateJWT(id int64, timeout time.Duration, secret []byte) (token string, err error) {
	claim := JWTClaims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(timeout).Unix(),
			Issuer:    Config.Issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	unsigned := jwt.NewWithClaims(getSigningMethod(Config.SigningAlgorithm), claim)
	token, err = unsigned.SignedString(secret)
	return
}

func GenerateAccessJWT(id int64) (string, error) {
	return generateJWT(id, Config.Timeout, accessSecret)
}

func GenerateRefreshJWT(id int64) (string, error) {
	return generateJWT(id, Config.MaxRefreshTime, refreshSecret)
}

// validateJWT validates whether jwt is valid.
// If so, we still have to check if user really logged in before.
func validateJWT(tokenString string, secret []byte) (int64, int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validation the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	claims := token.Claims.(jwt.MapClaims)
	if claims["id"] == nil || err != nil {
		return 0, 0, errs.ErrInvalidToken
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return 0, 0, errs.ErrInvalidToken
	}

	iat := claims["iat"].(float64)
	return int64(id), int64(iat), nil
}

func ValidateAccessJWT(token string) (int64, int64, error) {
	return validateJWT(token, accessSecret)
}

func ValidateRefreshJWT(token string) (int64, int64, error) {
	return validateJWT(token, refreshSecret)
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
