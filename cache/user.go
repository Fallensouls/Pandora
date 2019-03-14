package cache

import (
	. "github.com/go-pandora/core/conf"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	Login     = 1
	Logout    = 0
	PrefixJWT = "jwt"
)

func addPrefix(key string) string {
	return strings.Join([]string{PrefixJWT}, key)
}

// You can use bitmap of cache to keep a record of each user's login status.
// Here we don't use it.
func SetStatusLogin(id int64) error {
	return client.SetBit("online_user", id, Login).Err()
}

func SetStatusLogout(id int64) error {
	return client.SetBit("online_user", id, Logout).Err()
}

// RevokeJWT sets a deadline of user's jwt.
// All jwt issued before the new deadline will be revoked.
func RevokeJWT(id string) error {
	key := addPrefix(id)
	return client.Set(key, time.Now().Unix(), Config.Timeout).Err()
}

// IsJWTRevoked checks if user's jwt is revoked.
func IsJWTRevoked(id string, timestamp int64) (interface{}, bool) {
	key := addPrefix(id)
	unixTime, err := client.Get(key).Result()
	if err != nil {
		return id, false
	}
	log.Println(unixTime)
	loginTime, _ := strconv.ParseInt(unixTime, 10, 64)
	if timestamp < loginTime-3 {
		return nil, true
	}
	return id, false
}
