package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	Org, err := DBConnection.Query("select ID from 'organization' where (houseID = '" + strconv.Itoa(idH) + "')")
	if err != nil {
		return nil, err
	}

	var orgBuf []orgID

	for Org.Next() {
		t := orgID{}
		err := Org.Scan(&t.ID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		orgBuf = append(orgBuf, t)
	}

	//Result
	result := []orgInfo{}
	for _, v := range orgBuf {
		org, err := GetOrganizationFromDB(DBConnection, v.ID)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		result = append(result, org)
	}

	return result, nil
}
