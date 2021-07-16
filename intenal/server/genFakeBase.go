package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"simple2gis/intenal/sqlite"
	"strconv"
)

//Заполнение базы фейковыми данными (можно было написать отдельную тулзню, но решил в качестве метода, т.к. увлекся ¯\_(ツ)_/¯)

func genString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZёйцукенгшщзхъфывапролджэячсмитьбюЁЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЯЧСМИТЬБЮ1234567890")
	b := make([]rune, rand.Intn(n)+1)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func genInt(n int) int {
	res := 0
	for i := 0; i < n; i++ {
		res = (res * 10) + rand.Intn(9)
	}
	return res
}

func genNumber() []string {
	var rez []string
	for i := 0; i < rand.Intn(10); i++ {
		rez = append(rez, strconv.Itoa(genInt(11)))
	}
	return rez
}

func genCat() []string {
	var rez []string
	for i := 0; i < rand.Intn(5); i++ {
		rez = append(rez, genString(7))
	}
	return rez
}

//Да, просто, запускаем number раз добавление организации
//Какой-то хитрой логики тут нет, так что генерируемые данные не шипко обладают смыслом
func (s *HTTPServer) genFakeBase(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.URL)

	number, err := strconv.Atoi(r.Form.Get("number"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	for i := 0; i < number; i++ {
		sqlite.SetHouse(w, genString(10), genString(20), fmt.Sprintf("%v.%v", genInt(2), genInt(8)),
			fmt.Sprintf("%v.%v", genInt(2), genInt(8)),
			genNumber(), genCat())
	}
}
