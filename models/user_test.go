package models

import (
	"github.com/go-pandora/core/errs"
	"github.com/go-pandora/core/util/csvutil"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	cleanData()
	m.Run()
	cleanData()
}

func cleanData() {
	_, err := engine.Exec("truncate `users` restart identity")
	if err != nil {
		log.Println(err)
	}
}

func TestUser_AddUser(t *testing.T) {
	records, err := csvutil.GetTestData("../test/user_test.csv")
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

func TestUser_UpdateUserProfile(t *testing.T) {
	email := "pandora@gamil.com"
	user := User{Username: "pandora", Password: "pandora", Age: 20, Email: &email, Description: "I am a programmer.",
		Gender: Male}

	if _, err := engine.Table("users").Insert(&user); err != nil {
		t.Fatal(err)
	}
	user.Gender = Female
	user.Description = "I am a farmer."
	email = "pandora@qq.com"
	user.Email = &email
	if err := user.UpdateUserProfile(user.Id); err != nil {
		t.Error(err)
	}

	var updateUser User
	if _, err := engine.Table("users").ID(user.Id).Get(&updateUser); err != nil {
		t.Error(err)
	}

	assert := assert.New(t)
	assert.Equal(user.Gender, updateUser.Gender)
	assert.Equal(user.Description, updateUser.Description)
	assert.NotEqual(user.Email, updateUser.Email) // email address won't be changed.
}

func TestUser_Login(t *testing.T) {
	email1 := "Pandora3@gmail.com"
	email2 := "Pandora4@gmail.com"
	cellphone1 := "13343643535"
	cellphone2 := "13347674645"

	users := []*User{
		{Username: "Pandora1", Password: "Pandora&", Email: &email1, Status: Inactive},
		{Username: "Pandora2", Password: "pandora^-", Cellphone: &cellphone1, Status: Normal},
		{Username: "Pandora1", Password: "Pandora&", Email: &email2, Status: Restricted},
		{Username: "Pandora2", Password: "pandora^-", Cellphone: &cellphone2, Status: Banned},
	}

	for _, user := range users {
		hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(hash)
		if _, err := engine.Table("users").Insert(user); err != nil {
			t.Fatal(err)
		}
	}

	loginUsers := []User{
		{Email: &email1, Password: "Pandora&"},
		{Email: &email2, Password: "Pandora&"},
		{Cellphone: &cellphone1, Password: "pandora^-"},
		{Cellphone: &cellphone2, Password: "pandora^-"},
	}

	var e []error
	for _, loginUser := range loginUsers {
		err := loginUser.Login()
		e = append(e, err)
	}

	assert := assert.New(t)
	assert.Equal(errs.ErrUserInactive, e[0])
	assert.Nil(e[1])
	assert.Nil(e[2])
	assert.Equal(errs.ErrUserBanned, e[3])

	incorrectUser := User{Email: &email2, Password: "aaaaaaaa"}
	err := incorrectUser.Login()
	assert.Equal(errs.ErrWrongPassword, err)
}
