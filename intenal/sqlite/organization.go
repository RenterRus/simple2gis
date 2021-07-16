package sqlite

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

//Запихивает запрос в очередь
func GetOrganization(w http.ResponseWriter, id string) {
	waitingOrg := make(chan int)
	queue <- list{
		Organization: Organization{
			id: id,
			f: func(info orgInfo, err error) {
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
				} else {

					resp, errm := json.Marshal(info)
					if errm != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(errm.Error()))
						waitingOrg <- 1
						return
					}
					w.WriteHeader(http.StatusOK)
					w.Write(resp)

				}

				waitingOrg <- 1
			},
		},
	}
	<-waitingOrg
}

//Выполняет запрос, когда дойдет очередь
func getOrgByID(DBConnection *sql.DB, id string) (GetOrgInfo, error) {
	org, err := DBConnection.Query("select ID, Name, CatIDs, houseID from 'organization' where (ID = '" + id + "')")
	if err != nil {
		return GetOrgInfo{}, err
	}

	var result GetOrgInfo

	for org.Next() {
		err := org.Scan(&result.Id, &result.Name, &result.CatIDs, &result.HouseID)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return result, nil
}

func getHouseByOrgInfo(DBConnection *sql.DB, orgBuf GetOrgInfo) (GetHouseInfo, error) {
	houseReq, err := DBConnection.Query("select ID, Address, Geo from 'house' where (ID = '" + orgBuf.HouseID + "')")
	if err != nil {
		return GetHouseInfo{}, err
	}

	var result GetHouseInfo

	for houseReq.Next() {
		err := houseReq.Scan(&result.Id, &result.Address, &result.Geo)
		if err != nil {
			log.Println(err)
			continue
		}
	}
	return result, nil
}

func getNumberByOrgId(DBConnection *sql.DB, orgID string) ([]string, error) {
	numberReq, err := DBConnection.Query("select Number from 'number' where (organization_id = '" + orgID + "')")
	if err != nil {
		return nil, err
	}

	var result []string

	for numberReq.Next() {
		t := GetNumberInfo{}
		err := numberReq.Scan(&t.Telnum)
		if err != nil {
			log.Println(err)
			continue
		}
		result = append(result, t.Telnum)
	}
	return result, nil
}

func getIDCategorybyCategory(DBConnection *sql.DB, catID string) []string {
	var result []string

	for _, v := range strings.Split(catID, "/") {
		categoryReq := DBConnection.QueryRow("select category from 'category' where (ID = " + v + ")")
		t := GetCategoryInfo{}
		categoryReq.Scan(&t.Category)
		result = append(result, t.Category)
	}
	return result
}

func GetOrganizationFromDB(DBConnection *sql.DB, id string) (orgInfo, error) {
	//Org
	orgBuf, err := getOrgByID(DBConnection, id)
	if err != nil {
		return orgInfo{}, err
	}

	//House
	houseBuf, err := getHouseByOrgInfo(DBConnection, orgBuf)
	if err != nil {
		return orgInfo{}, err
	}

	//Number
	numberBuf, err := getNumberByOrgId(DBConnection, orgBuf.Id)
	if err != nil {
		return orgInfo{}, err
	}

	//Category
	categoryBuf := getIDCategorybyCategory(DBConnection, orgBuf.CatIDs)

	//Response
	var response orgInfo

	if orgBuf.Id != "" {
		response.Addr = houseBuf.Address
		houseGeo := strings.Split(houseBuf.Geo, "/")
		response.GeoX = houseGeo[0]
		response.GeoY = houseGeo[1]
		response.Name = orgBuf.Name
		response.Numbers = numberBuf
		response.Categorys = strings.Join(categoryBuf, "/")
	}

	return response, nil
}
