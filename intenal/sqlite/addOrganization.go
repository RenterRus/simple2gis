package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//Формат ответа
type addInfo struct {
	Message string `json:"Message"`
}

//Запихивает запрос в канал и начинает ждать ответа (т.к. каждый запрос через http генерирует новую горутину, то это не блокирует метод сервера)
func SetHouse(w http.ResponseWriter, name, houseAddr, houseGeoX, houseGeoY string, numbers, category []string) {
	waitingAdd := make(chan int)
	queue <- list{
		AddHouse: AddHouse{
			name:      name,
			numbers:   numbers,
			houseAddr: houseAddr,
			houseGeoX: houseGeoX,
			houseGeoY: houseGeoY,
			category:  category,
			f: func(info addInfo, err error) {
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
				} else {

					resp, errm := json.Marshal(info)
					if errm != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(errm.Error()))
						waitingAdd <- 1
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write(resp)

				}

				waitingAdd <- 1
			},
		},
	}
	<-waitingAdd
}

//Все что ниже - заполняет БД данными, вызывается ProcAddHouse (оркестратор данной функции) через оркестратор очереди (worker)
type housesID struct {
	Id int
}

func GetHouseID(DBConnection *sql.DB, address, geo string) (int, error) {
	result, err := DBConnection.Query("select ID from house where (Address = '" + address + "' AND Geo = '" + geo + "')")
	if err != nil {
		return -2, err
	}

	var houses []housesID

	for result.Next() {
		t := housesID{}
		err := result.Scan(&t.Id)
		if err != nil {
			log.Println(err)
			continue
		}
		houses = append(houses, t)
	}
	if len(houses) > 0 {
		return houses[0].Id, nil
	}
	id, erri := AddHouseToDB(DBConnection, address, geo)
	if erri != nil {
		return -3, erri
	}
	return id, nil
}

func AddHouseToDB(DBConnection *sql.DB, address, geo string) (int, error) {
	res, err := DBConnection.Exec("insert into 'house' (Address, Geo) VALUES ('" + address + "', '" + geo + "')")
	if err != nil {
		return -1, err
	}
	id, erri := res.LastInsertId()
	if erri != nil {
		return -1, err
	}
	return int(id), nil
}

func AddNumbersToDB(DBConnection *sql.DB, numbers []string, organizationID int) error {
	for _, v := range numbers {
		_, err := DBConnection.Exec("insert into 'number' (Number, organization_id) VALUES ('" + v + "', '" + strconv.Itoa(organizationID) + "')")
		if err != nil {
			return err
		}
	}

	return nil
}

func GetRootCategoryFromDB(DBConnection *sql.DB, cat string, rootID int) (int, error) {
	result, err := DBConnection.Query("select ID from 'category' where (category = '" + cat + "' and RootID = '" + strconv.Itoa(rootID) + "')")
	if err != nil {
		return -2, err
	}

	var houses []string

	for result.Next() {
		t := housesID{}
		err := result.Scan(&t.Id)
		if err != nil {
			log.Println(err)
			continue
		}
		houses = append(houses, strconv.Itoa(t.Id))
	}
	if len(houses) > 0 {
		return strconv.Atoi(houses[0])
	}
	res, err := DBConnection.Exec("insert into 'category' (category, RootID) VALUES ('" + cat + "', '" + strconv.Itoa(rootID) + "')")
	if err != nil {
		return -3, err
	}
	resID, err := res.LastInsertId()
	if err != nil {
		return -4, err
	}

	return int(resID), nil
}

func AddCategoryToDB(DBConnection *sql.DB, cats []string) ([]string, error) {
	id := 0
	var resId []string
	var err error
	for _, v := range cats {
		id, err = GetRootCategoryFromDB(DBConnection, v, id)
		if err != nil {
			return nil, err
		}
		resId = append(resId, strconv.Itoa(id))
	}
	return resId, nil
}

func AddOrganizationToDB(DBConnection *sql.DB, name, catIDs string, houseID int) (int, error) {
	res, err := DBConnection.Exec("insert into 'organization' (Name, CatIDs, houseID) VALUES ('" + name + "', '" + catIDs + "', '" + strconv.Itoa(houseID) + "')")
	if err != nil {
		return -1, err
	}

	result, err := res.LastInsertId()
	return int(result), err
}

type doubleID struct {
	ID string `json:"id"`
}

type houseID struct {
	Address string
	Geo     string
}

func checkDouble(DBConnection *sql.DB, name, houseAddr, houseGeoX, houseGeoY string) bool {
	result, err := DBConnection.Query("select houseID from 'organization' where (Name = '" + name + "')")
	if err != nil {
		log.Println(err.Error())
	}
	var orgs []doubleID

	for result.Next() {
		t := doubleID{}
		err := result.Scan(&t.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		orgs = append(orgs, t)
	}
	answer := false
	for _, v := range orgs {
		result, _ = DBConnection.Query("select Address, Geo from 'house' where (ID = '" + v.ID + "')")

		for result.Next() {
			t := houseID{}
			err := result.Scan(&t.Address, &t.Geo)
			if err != nil {
				log.Println(err)
				continue
			}
			if t.Address == houseAddr && t.Geo == houseGeoX+"/"+houseGeoY {
				answer = true
			}
		}
	}

	return answer
}

func ProcAddHouse(DBConnection *sql.DB, name, houseAddr, houseGeoX, houseGeoY string, numbers, category []string) (string, error) {
	if !checkDouble(DBConnection, name, houseAddr, houseGeoX, houseGeoY) {
		hId, err := GetHouseID(DBConnection, houseAddr, houseGeoX+"/"+houseGeoY)
		if err != nil {
			return err.Error(), err
		}

		catIDs, err := AddCategoryToDB(DBConnection, category)
		if err != nil {
			return err.Error(), err
		}

		orgIDres, err := AddOrganizationToDB(DBConnection, name, strings.Join(catIDs, "/"), hId)
		if err != nil {
			return err.Error(), err
		}

		err = AddNumbersToDB(DBConnection, numbers, orgIDres)
		if err != nil {
			return err.Error(), err
		}

		return fmt.Sprintf("Completed successfully: %v", orgIDres), nil
	} else {
		return "Already exists", nil
	}
}
