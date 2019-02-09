package redis

import (
	"strconv"
	"time"
)

const (
	Login  = 1
	Logout = 0
)

// You can use bitmap of redis to keep a record of each user's login status.

func SetStatusLogin(id int64) error {
	return client.SetBit("LoginStatus", id, Login).Err()
}

func SetStatusLogout(id int64) error {
	return client.SetBit("LoginStatus", id, Logout).Err()
}

// SetNBFTime resets NotBefore time of user's jwt.
// All jwt issued before the new NotBefore time will be rejected.
func SetNBFTime(id int64) error {
	return client.HSet("nbf", strconv.FormatInt(id, 10), time.Now().Unix()).Err()
}

// CheckJWTInBlacklist checks if user's jwt is in blacklist.
func CheckJWTInBlacklist(id int64, timestamp int64) (bool, error) {
	unixTime, err := client.HGet("nbf", strconv.FormatInt(id, 10)).Result()
	if err != nil {
		return false, err
	}
	loginTime, _ := strconv.ParseInt(unixTime, 10, 64)
	if timestamp < loginTime-3 {
		return false, nil
	}
	return true, nil
}
