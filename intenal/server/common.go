package server

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

//Склад непрофильных функций и переменных
var regexpTelnum = regexp.MustCompile("^[0-9]{10,11}")

func CheckID(id string) (string, error) {
	if _, err := strconv.Atoi(id); err != nil || id[0] == '0' {
		return "", errors.New("ID is invalid")
	}
	return id, nil
}

func (s *HTTPServer) agregator(form url.Values) map[string]string {
	httpArgs := map[string]string{}
	for k, v := range form {
		httpArgs[k] = strings.Join(v, "")
	}
	return httpArgs
}

func checkTelnum(telnum string) (string, error) {
	if regexpTelnum.MatchString(telnum) {
		return telnum, nil
	}

	return "", errors.New("invalid telnum: " + telnum)
}
