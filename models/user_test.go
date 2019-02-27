package models

import (
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

func TestUser_GetUser(t *testing.T) {
	email1 := "Pandora1@gmail.com"
	email2 := "Pandora2@gmail.com"
	cellphone1 := "13345643535"
	cellphone2 := "13345674645"

	users := []*User{
		{Username: "Pandora1", Password: "Pandora&", Email: &email1, Status: Inactive},
		{Username: "Pandora2", Password: "pandora^-", Cellphone: &cellphone1, Status: Normal},
		{Username: "Pandora1", Password: "Pandora&", Email: &email2, Status: Restricted},
		{Username: "Pandora2", Password: "pandora^-", Cellphone: &cellphone2, Status: Banned},
	}

	for _, user := range users {
		if _, err := engine.Table("users").Insert(user); err != nil {
			t.Fatal(err)
		}
	}

	assert := assert.New(t)

	assert.Equal(errs.ErrUserInactive, new(User).GetUser(users[0].Id))
	assert.Equal(errs.ErrUserRestricted, new(User).GetUser(users[2].Id))
	assert.Equal(errs.ErrUserBanned, new(User).GetUser(users[3].Id))

	assert.Nil(new(User).GetUser(users[1].Id))
	assert.Equal(errs.ErrUserNotFound, new(User).GetUser(0))

}
