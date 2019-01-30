package validation

import (
	"errors"
	"fmt"
	"regexp"
)

func ValidateUsername(username string) error {
	re, _ := regexp.Compile(`^\w{4,16}$`)
	if !re.MatchString(username){
		return errors.New("username is invalid")
	}
	return nil
}

func ValidatePassword(password string) error {
	//re, err := regexp.Compile(`^(?!\d+$)(?![a-zA-Z]+$)(?![_~@#$^]+$)[\w~@#$^]{8,16}$`)
	re, err := regexp.Compile(`^[\w~@#$^]{8,16}$`)
	if err != nil{
		fmt.Println(err)
	}
	if !re.MatchString(password){
		return errors.New("password is invalid")
	}
	return nil
}

func ValidateEmail(email string) error {
	re, _ := regexp.Compile(`^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z]{2,6}$`)
	if !re.MatchString(email){
		return errors.New("email address is invalid")
	}
	return nil
}

func ValidateCellphone(tel string) error {
	re, _ := regexp.Compile(`^1([3[0-9]|4[579]|5[0-3,5-9]|6[6]|7[0135678]|8[0-9]|9[89])\d{8}$`)
	if !re.MatchString(tel){
		return errors.New("cellphone number is invalid")
	}
	return nil
}

