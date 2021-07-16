package sqlite

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

//Запихивает запрос в очередь
func GetHouse(w http.ResponseWriter, addr, x, y string) {
	waitingHouse := make(chan int)
	queue <- list{
		House: House{
			addr: addr,
			x:    x,
			y:    y,
			f: func(info []orgInfo, err error) {
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
				} else {

					resp, errm := json.Marshal(info)
					if errm != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(errm.Error()))
						waitingHouse <- 1
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write(resp)

				}
				waitingHouse <- 1
			},
		},
	}
	<-waitingHouse
}

//Выполняет запрос, когда дойдет очередь
type orgID struct {
	ID string `json:"ID"`
}

func GetOrgByHouse(DBConnection *sql.DB, Addr, Geo string) ([]orgInfo, error) {
	idH, err := GetHouseID(DBConnection, Addr, Geo)
	if err != nil {
		return nil, err
	}

	//Org
	org, err := DBConnection.Query("select ID from 'organization' where (houseID = '" + strconv.Itoa(idH) + "')")
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
