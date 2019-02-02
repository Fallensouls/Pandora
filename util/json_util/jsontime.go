package json_util

import "time"

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(j).Format("2006-01-02 15:04:05") + `"`), nil
}

func (j *JsonTime) GetJsonTime() {
	*j = JsonTime(time.Unix(time.Now().Unix(), 0))
}
