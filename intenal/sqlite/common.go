package sqlite

import "database/sql"

type House struct {
	addr string
	x    string
	y    string
	f    func(info []orgInfo, err error)
}

type AddHouse struct {
	name      string
	numbers   []string
	houseAddr string
	houseGeoX string
	houseGeoY string
	category  []string
	f         func(info addInfo, err error)
}

type Category struct {
	cat []string
	f   func(info []orgInfo, err error)
}

type Another struct {
	addr string
}

type Organization struct {
	id string
	f  func(info orgInfo, err error)
}

type list struct {
	House        House
	Organization Organization
	Another      Another
	Category     Category
	AddHouse     AddHouse
}

//Склад общих переменных и функций
var (
	queue    = make(chan list)
	DBClient *Database
)

type Database struct {
	DBName       string
	DBConnection *sql.DB
}

type GetOrgInfo struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	CatIDs  string `json:"catIDs"`
	HouseID string `json:"houseID"`
}

type GetHouseInfo struct {
	Id      string `json:"id"`
	Address string `json:"Address"`
	Geo     string `json:"Geo"`
}

type GetNumberInfo struct {
	Telnum string `json:"Number"`
}

type GetCategoryInfo struct {
	ID       string `json:"ID"`
	Category string `json:"Category"`
}

type orgInfo struct {
	Name      string   `json:"Name"`
	Addr      string   `json:"Address"`
	GeoX      string   `json:"X"`
	GeoY      string   `json:"Y"`
	Numbers   []string `json:"Numbers"`
	Categorys string   `json:"Category"`
}
