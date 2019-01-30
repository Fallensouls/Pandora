package date

import "time"

func GetStandardTime() time.Time {
	return time.Unix(time.Now().Unix(), 0)
}