package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

//Получение коннекта к DB (SQLite)
func Initial(dbname string) *Database {
	d := new(Database)
	d.DBName = dbname
	var err error
	d.DBConnection, err = sql.Open("sqlite3", d.DBName)
	if err != nil {
		panic(err.Error())
	}
	return d
}

//Красивое завершение
func (d *Database) DisableConnect() {
	d.DBConnection.Close()
}

//Оркестратор очередей (т.к. SQLite очень капризная к многопоточной работе, то все запросы выстраиваются в последовательную очередь)
//Запускается в main, завершается там же
func DBProc() {
	for {
		t := <-queue
		if t.Another.addr != "" {
			log.Println("Another")
		}
		if t.House.addr != "" {
			log.Println("House")
			t.House.f(GetOrgByHouse(DBClient.DBConnection, t.House.addr, t.House.x+"/"+t.House.y))
		}
		if t.Organization.id != "" {
			log.Println("Organization")

			t.Organization.f(GetOrganizationFromDB(DBClient.DBConnection, t.Organization.id))
		}
		if t.Category.cat != nil {
			log.Println("Category")

			t.Category.f(GetOrgByCategory(DBClient.DBConnection, t.Category.cat))
		}
		if t.AddHouse.name != "" {
			log.Println("AddHouse")
			res, err := ProcAddHouse(DBClient.DBConnection, t.AddHouse.name, t.AddHouse.houseAddr, t.AddHouse.houseGeoX,
				t.AddHouse.houseGeoY, t.AddHouse.numbers, t.AddHouse.category)
			t.AddHouse.f(addInfo{
				Message: res,
			}, err)
		}
	}
}
