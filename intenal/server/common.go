package server

import (
	"errors"
	"regexp"
	"strconv"
)

//Склад непрофильных функций и переменных
var regexpTelnum = regexp.MustCompile("^[0-9]{10,11}")

func CheckID(id string) (string, error) {
	if _, err := strconv.Atoi(id); err != nil || id[0] == '0' {
		return "", errors.New("ID is invalid")
	}
	return id, nil
}

func checkTelnum(telnum string) (string, error) {
	if regexpTelnum.MatchString(telnum) {
		return telnum, nil
	}

	return "", errors.New("invalid telnum: " + telnum)
}
