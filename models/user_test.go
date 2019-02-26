package models

import (
	"fmt"
	"github.com/Fallensouls/Pandora/errs"
	"github.com/Fallensouls/Pandora/util/csvutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_AddUser(t *testing.T) {
	records, err := csvutil.GetTestData("test/user_test.csv")
	if err != nil {
		t.Fatal("no test data")
	}

	var (
		data []User
		e    []error
	)

	for _, record := range records {
		fmt.Println(record)
		switch {
		case record[2] == "" && record[3] == "":
			data = append(data, User{Username: record[0], Password: record[1], Email: nil, Cellphone: nil})
		case record[2] == "":
			data = append(data, User{Username: record[0], Password: record[1], Email: nil, Cellphone: &record[3]})
		case record[3] == "":
			data = append(data, User{Username: record[0], Password: record[1], Email: &record[2], Cellphone: nil})
		default:
			data = append(data, User{Username: record[0], Password: record[1], Email: &record[2], Cellphone: &record[3]})
		}
		e = append(e, errs.ErrMap[record[4]])
	}
	data = data[1:]
	e = e[1:]

	assert := assert.New(t)
	for i, user := range data {
		err := user.AddUser()
		assert.Equal(e[i], err)
	}
}
