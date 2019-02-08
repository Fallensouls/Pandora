package redis

import "log"

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

// CheckStatus check if user have logged in.
func CheckLoginStatus(id int64) (bool, error) {
	status, err := client.GetBit("login_status", id).Result()
	log.Println(status)
	if status == Login {
		return true, err
	}
	return false, err
}
