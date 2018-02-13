package main

import (
	"database/sql"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

/*Available Actions:
Get - get entire table
GetActive - get rows which are not submitted in previous daily reports
Update - update a row
Append - append a new row to the end
Delete - delete a row
Replace - flush table and input new data into table
ReplaceActive - flush active record and append new data into table
UpdateStatus - change status of order/purchase from active "0" to processed "1"
*/

//unit
func unittableGet() (units []string) {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-unittableGet-1")
	defer db.Close()
	rows, err := db.Query("SELECT * FROM unittable")
	checkError(err, "model-unittableGet-2")

	var unit string
	for rows.Next() {
		err = rows.Scan(&unit)
		checkError(err, "model-unittableGet-3")
		units = append(units, unit)
	}
	return
}

func unittableReplace(units []string) bool {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-unittableReplace-1")
	defer db.Close()
	_, err = db.Exec("delete from unittable")
	checkError(err, "model-unittableReplace-2")

	stmt, err := db.Prepare("insert into unittable(unit) values(?)")
	checkError(err, "model-unittableReplace-3")
	defer stmt.Close()
	for _, unit := range units {
		_, err = stmt.Exec(unit)
		checkError(err, "model-unittableReplace-4")
	}
	return true
}

//item
func itemtableGet() (items []string, units []string, unitprices []float64, notes []string) {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-itemtableGet-1")
	defer db.Close()
	rows, err := db.Query("SELECT * FROM itemtable")
	checkError(err, "model-itemtableGet-2")

	var unitprice float64
	var item, unit, note string
	for rows.Next() {
		err = rows.Scan(&item, &unit, &unitprice, &note)
		checkError(err, "model-itemtableGet-3")
		items = append(items, item)
		units = append(units, unit)
		unitprices = append(unitprices, unitprice)
		notes = append(notes, note)
	}
	return
}

func itemtableReplace(items []string, units []string, unitprices []float64, notes []string) bool {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-itemtableReplace-1")
	defer db.Close()
	_, err = db.Exec("delete from itemtable")
	checkError(err, "model-itemtableReplace-2")

	stmt, err := db.Prepare("insert into itemtable(item,unit,unitprice,notes) values(?,?,?,?)")
	checkError(err, "model-itemtableReplace-3")
	defer stmt.Close()
	for i, _ := range items {
		_, err = stmt.Exec(items[i], units[i], unitprices[i], notes[i])
		checkError(err, "model-itemtableReplace-4")
	}
	return true
}

//bldg
func bldgtableGet() (regdates []int, bldgs []string, addrs []string, zips []int, notes []string) {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-bldgtableGet-1")
	defer db.Close()
	rows, err := db.Query("SELECT * FROM bldgtable")
	checkError(err, "model-bldgtableGet-2")

	var regdate, zip int
	var bldg, addr, note string
	for rows.Next() {
		err = rows.Scan(&regdate, &bldg, &addr, &zip, &note)
		checkError(err, "model-bldgtableGet-3")
		regdates = append(regdates, regdate)
		bldgs = append(bldgs, bldg)
		addrs = append(addrs, addr)
		zips = append(zips, zip)
		notes = append(notes, note)
	}
	return
}

func bldgtableReplace(regdates []int, bldgs []string, addrs []string, zips []int, notes []string) bool {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-bldgtableReplace-1")
	defer db.Close()
	_, err = db.Exec("delete from bldgtable")
	checkError(err, "model-bldgtableReplace-2")

	stmt, err := db.Prepare("insert into bldgtable(regdate, bldg, addr, zip, notes) values(?,?,?,?,?)")
	checkError(err, "model-bldgtableReplace-3")
	defer stmt.Close()
	for i, _ := range regdates {
		_, err = stmt.Exec(regdates[i], bldgs[i], addrs[i], zips[i], notes[i])
		checkError(err, "model-bldgtableReplace-4")
	}
	return true
}

//customer
func customertableGet() (ids []int, regdates []int, nicknames []string, phones []int64, bldgs []string, rooms []string, notes []string) {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-customertableGet-1")
	defer db.Close()
	rows, err := db.Query("SELECT * FROM customertable")
	checkError(err, "model-customertableGet-2")

	var id, regdate int
	var phone int64
	var nickname, bldg, room, note string
	for rows.Next() {
		err = rows.Scan(&id, &regdate, &nickname, &phone, &bldg, &room, &note)
		checkError(err, "model-customertableGet-3")
		ids = append(ids, id)
		regdates = append(regdates, regdate)
		nicknames = append(nicknames, nickname)
		phones = append(phones, phone)
		bldgs = append(bldgs, bldg)
		rooms = append(rooms, room)
		notes = append(notes, note)
	}
	return
}

func customertableAppend(regdate int, nickname string, phone int64, bldg string, room string, note string) bool {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-customertableAppend-1")
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO customertable(regdate, nickname, phone, bldg, room, notes) values(?,?,?,?,?,?)")
	checkError(err, "model-customertableAppend-2")

	_, err = stmt.Exec(regdate, nickname, phone, bldg, room, note)
	checkError(err, "model-customertableAppend-3")
	return true
}

func customertableUpdate(id int, regdate int, nickname string, phone int64, bldg string, room string, note string) bool {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-customertableUpdate-1")
	defer db.Close()
	stmt, err := db.Prepare("update customertable set regdate=?, nickname=?, phone=?, bldg=?, room=?, notes=?  where id=?")
	checkError(err, "model-customertableUpdate-2")

	_, err = stmt.Exec(regdate, nickname, phone, bldg, room, note, id)
	checkError(err, "model-customertableUpdate-3")
	return true
}

func customertableDelete(id int) bool {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-customertableDelete-1")
	defer db.Close()
	stmt, err := db.Prepare("delete from customertable where id=?")
	checkError(err, "model-customertableDelete-2")

	_, err = stmt.Exec(id)
	checkError(err, "model-customertableDelete-3")
	return true
}

//order
func ordertableGetActive() (ids []int, nicknames []string, orderlists []string) {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-ordertableGetActive-1")
	defer db.Close()
	rows, err := db.Query("SELECT * FROM ordertable WHERE status=0")
	checkError(err, "model-ordertableGetActive-2")

	var id, tempdate, status int
	var nickname, orderlist string
	for rows.Next() {
		err = rows.Scan(&id, &nickname, &tempdate, &orderlist, &status)
		checkError(err, "model-ordertableGetActive-3")
		ids = append(ids, id)
		nicknames = append(nicknames, nickname)
		orderlists = append(orderlists, orderlist)
	}
	return
}

func ordertableGetActivePerNickname(nickname string) (ids []int, dates []int) {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-ordertableGetActiveIdsPerNickname-1")
	defer db.Close()
	rows, err := db.Query("SELECT id, orderdate FROM ordertable WHERE status=0 AND nickname='" + nickname + "'")
	checkError(err, "model-ordertableGetActiveIdsPerNickname-2")

	var id, date int
	for rows.Next() {
		err = rows.Scan(&id, &date)
		checkError(err, "model-ordertableGetActiveIdsPerNickname-3")
		ids = append(ids, id)
		dates = append(dates, date)
	}
	return
}

func ordertableAppend(nickname string, orderdate int, orderlist string) bool {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-ordertableAppend-1")
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO ordertable(nickname, orderdate, orderlist, status) values(?,?,?,?)")
	checkError(err, "model-ordertableAppend-2")

	_, err = stmt.Exec(nickname, orderdate, orderlist, 0)
	checkError(err, "model-ordertableAppend-3")
	return true
}

func ordertableDelete(ids []int) bool {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-ordertableDelete-1")
	defer db.Close()

	idStrs := make([]string, 0)
	for _, val := range ids {
		idStrs = append(idStrs, strconv.FormatInt(int64(val), 10))
	}
	idStr := strings.Join(idStrs, ",")

	_, err = db.Exec("delete from ordertable where id in (" + idStr + ")")
	checkError(err, "model-ordertableDelete-2")
	return true
}

//purchase
func purchasetableGetActive() (ids []int, dates []int, items []string, units []string, amounts []float64, prices []float64) {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-purchasetableGetActive-1")
	defer db.Close()
	rows, err := db.Query("SELECT * FROM purchasetable WHERE status=0")
	checkError(err, "model-purchasetableGetActive-2")

	var id, date, status int
	var item, unit string
	var amount, price float64
	for rows.Next() {
		err = rows.Scan(&id, &date, &item, &unit, &amount, &price, &status)
		checkError(err, "model-purchasetableGetActive-3")
		ids = append(ids, id)
		dates = append(dates, date)
		items = append(items, item)
		units = append(units, unit)
		amounts = append(amounts, amount)
		prices = append(prices, price)
	}
	return
}
func purchasetableAppend(purchasedate int, item string, unit string, amount float64, price float64) bool {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-purchasetableAppend-1")
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO purchasetable(purchasedate, item, unit, amount, price, status) values(?,?,?,?,?,?)")
	checkError(err, "model-purchasetableAppend-2")

	_, err = stmt.Exec(purchasedate, item, unit, amount, price, 0)
	checkError(err, "model-purchasetableAppend-3")
	return true
}

func purchasetableReplaceActive(purchasedate int, items []string, units []string, amounts []float64, prices []float64) bool {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-purchasetableReplaceActive-1")
	defer db.Close()

	_, err = db.Exec("delete from purchasetable where status=0")
	checkError(err, "model-purchasetableReplaceActive-2")

	stmt, err := db.Prepare("insert into purchasetable(purchasedate, item, unit, amount, price, status) values(?,?,?,?,?,?)")
	checkError(err, "model-purchasetableReplaceActive-3")
	defer stmt.Close()
	for i, _ := range items {
		_, err = stmt.Exec(purchasedate, items[i], units[i], amounts[i], prices[i], 0)
		checkError(err, "model-purchasetableReplaceActive-4")
	}

	return true
}

//status
func orderpurchaseUpdateStatus() bool {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-orderpurchaseUpdateStatus-1")
	defer db.Close()

	var id int

	rows, err := db.Query("SELECT id FROM ordertable WHERE status=0")
	checkError(err, "model-orderpurchaseUpdateStatus-2")

	ids := make([]int, 0)
	for rows.Next() {
		err = rows.Scan(&id)
		checkError(err, "model-orderpurchaseUpdateStatus-4")
		ids = append(ids, id)
	}

	stmt, err := db.Prepare("update ordertable set status=?  where id=?")
	checkError(err, "model-orderpurchaseUpdateStatus-3")
	for _, val := range ids {
		_, err = stmt.Exec(1, val)
		checkError(err, "model-orderpurchaseUpdateStatus-5")
	}

	rows, err = db.Query("SELECT id FROM purchasetable WHERE status=0")
	checkError(err, "model-orderpurchaseUpdateStatus-6")

	ids = nil
	for rows.Next() {
		err = rows.Scan(&id)
		checkError(err, "model-orderpurchaseUpdateStatus-8")
		ids = append(ids, id)
	}

	stmt, err = db.Prepare("update purchasetable set status=?  where id=?")
	checkError(err, "model-orderpurchaseUpdateStatus-7")
	for _, val := range ids {
		_, err = stmt.Exec(1, val)
		checkError(err, "model-orderpurchaseUpdateStatus-9")
	}
	return true
}

func logTempReport(date int, purchases string, orders string) bool {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-logTempReport-1")
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO tempReport(date, purchaseList, orderList) values(?,?,?)")
	checkError(err, "model-logTempReport-2")

	_, err = stmt.Exec(date, purchases, orders)
	checkError(err, "model-logTempReport-3")
	return true
}

func extractTempReport() (date int, purchases string, orders string) {
	db, err := sql.Open("sqlite3", "./database/record.db")
	checkError(err, "model-extractTempReport-1")
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tempReport")
	checkError(err, "model-extractTempReport-2")
	for rows.Next() {
		err = rows.Scan(&date, &purchases, &orders)
		checkError(err, "model-extractTempReport-3")
	}
	_, err = db.Exec("delete from tempReport") //empty the temp report after being read
	checkError(err, "model-extractTempReport-4")
	return
}

func logReport(date int, purchases string, orders string) bool {
	db, err := sql.Open("sqlite3", "./database/report.db")
	checkError(err, "model-logReport-1")
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO reports(date, purchaseList, orderList) values(?,?,?)")
	checkError(err, "model-logReport-2")

	_, err = stmt.Exec(date, purchases, orders)
	checkError(err, "model-logReport-3")
	return true
}
