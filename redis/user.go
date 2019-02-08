package redis

import (
	"strconv"
	"time"
)

const (
	Login  = 1
	Logout = 0
)

func SetStatusLogin(id int64) error {
	return client.SetBit("login_status", id, Login).Err()
}

func SetStatusLogout(id int64) error {
	return client.SetBit("login_status", id, Logout).Err()
}

func SetLoginTime(id int64) error {
	return client.HSet("login_time", strconv.FormatInt(id, 10), time.Now().Unix()).Err()
}

// CheckJWTStatus checks if user's jwt becomes invalid.
func CheckJWTStatus(id int64, timestamp int64) (bool, error) {
	status, err := client.GetBit("login_status", id).Result()
	if status == Logout {
		return false, err
	}
	unixTime, err := client.HGet("login_time", strconv.FormatInt(id, 10)).Result()
	if err != nil {
		return false, err
	}
	loginTime, _ := strconv.ParseInt(unixTime, 10, 64)
	if timestamp < loginTime-3 {
		return false, nil
	}
	return true, nil
}
