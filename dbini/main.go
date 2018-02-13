package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	recordIni()
	reportIni()
}

func recordIni() {
	os.Remove("./record.db")
	db, err := sql.Open("sqlite3", "./record.db")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("create table unittable (unit text not null);")
	checkErr(err) //ea
	_, err = db.Exec("create table itemtable (item text not null, unit text not null, unitprice real not null, notes text);")
	checkErr(err) //apple, ea, 0.99, sweet
	_, err = db.Exec("create table bldgtable (regdate integer not null, bldg text not null, addr text not null, zip integer not null, notes text);")
	checkErr(err) //20160215, Harrison St. Senior House, 1000 Harrison St., 94555, good place
	_, err = db.Exec("create table customertable (id integer not null primary key autoincrement, regdate integer not null, nickname text not null, phone integer not null, bldg text not null, room text not null, notes text);")
	checkErr(err) //1, 20160215, John-HS-301, 5105698756, Harrison St. Senior House, 301A, good guy
	_, err = db.Exec("create table ordertable (id integer not null primary key autoincrement, nickname text not null, orderdate integer not null, orderlist text not null, status integer not null);")
	checkErr(err) //1, John-HS-301, 20160216, item ^ unit ^ amount ^ notes ^ unit2 ^ amount2 ^ price? ... , 0
	_, err = db.Exec("create table purchasetable (id integer not null primary key autoincrement, purchasedate integer not null, item text not null, unit text not null, amount real not null, price real not null, status integer not null);")
	checkErr(err) //1, 20160216, apple, ea, 12, 50, 0
	//status = 0 -> not submitted; status > 0 -> submitted
	_, err = db.Exec("create table tempReport (date integer not null, purchaseList text not null, orderList text not null)")
	checkErr(err)
}

func reportIni() {
	os.Remove("./report.db")
	db, err := sql.Open("sqlite3", "./report.db")
	checkErr(err)
	defer db.Close()

	_, err = db.Exec("create table reports (id integer not null primary key autoincrement, date integer not null, purchaseList text not null, orderList text not null)")
	checkErr(err)
}
