package sqlite

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

//Запихивает запрос в канал
func GetCategory(w http.ResponseWriter, cat []string) {
	waitingCat := make(chan int)
	queue <- list{
		Category: Category{
			cat: cat,
			f: func(info []orgInfo, err error) {
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
				} else {

					resp, errm := json.Marshal(info)
					if errm != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(errm.Error()))
						waitingCat <- 1
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write(resp)

				}

				waitingCat <- 1
			},
		},
	}
	<-waitingCat
}

//Все, что ниже - выполняет указания оркестратора очереди
type catId struct {
	ID string `json:"ID"`
}

func GetOrgByCategory(DBConnection *sql.DB, category []string) ([]orgInfo, error) {
	var catIDs []string
	for _, v := range category {
		result, err := DBConnection.Query("select ID from 'category' where (category = '" + v + "')")
		if err != nil {
			return nil, err
		}

		for result.Next() {
			t := catId{}
			err := result.Scan(&t.ID)
			if err != nil {
				log.Println(err)
				continue
			}
			catIDs = append(catIDs, t.ID)
		}
	}

	//Org
	org, err := DBConnection.Query("select ID from 'organization' where CatIDs like '%" + strings.Join(catIDs, "/") + "%'")
	if err != nil {
		return nil, err
	}

	var orgBuf []orgID

	for org.Next() {
		t := orgID{}
		err := org.Scan(&t.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		orgBuf = append(orgBuf, t)
	}

	//Result
	result := []orgInfo{}
	for _, v := range orgBuf {
		orgs, err := GetOrganizationFromDB(DBConnection, v.ID)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		result = append(result, orgs)
	}

	return result, nil
}
