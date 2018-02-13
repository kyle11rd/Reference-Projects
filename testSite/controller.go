package main

import (
	//"fmt"
	"sort"
	"strconv"
	"strings"
)

func getUnits() []string {
	units := unittableGet()
	return units
}

func updateUnits(units []string) bool {
	unitsMod := make([]string, 0)
	for _, unit := range units {
		if unit != "" {
			unitsMod = append(unitsMod, unit)
		}
	}
	_ = unittableReplace(unitsMod)
	return true
}

func getItems() (items []string, units []string, unitPrices []float64, notes []string) {
	items, units, unitPrices, notes = itemtableGet()
	return
}

func updateItems(itemsRaw []string, unitsRaw []string, unitPricesRaw []string, notesRaw []string) bool { //unitprice read from map as string
	items := make([]string, 0)
	units := make([]string, 0)
	unitPrices := make([]float64, 0)
	notes := make([]string, 0)

	for i, _ := range itemsRaw {
		if itemsRaw[i] != "" && unitsRaw[i] != "" && unitPricesRaw[i] != "" {
			items = append(items, itemsRaw[i])
			units = append(units, unitsRaw[i])
			tempPrice, err := strconv.ParseFloat(unitPricesRaw[i], 64)
			checkError(err, "controller-updateItems")
			unitPrices = append(unitPrices, tempPrice)
			notes = append(notes, notesRaw[i])
		}
	}
	_ = itemtableReplace(items, units, unitPrices, notes)
	return true
}

func getBldgs() (regdates []int, bldgs []string, addrs []string, zips []string, notes []string) {
	regdates, bldgs, addrs, zipsInt, notes := bldgtableGet()
	for _, val := range zipsInt { //convert zips to string as the form returns string
		zips = append(zips, strconv.Itoa(val))
	}
	return
}

func updateBldgs(bldgsRaw []string, addrsRaw []string, zipsRaw []string, notesRaw []string) bool { //sorted by bldg name
	bldgs := make([]string, 0)
	addrs := make([]string, 0)
	zips := make([]string, 0)
	notes := make([]string, 0)
	for i, _ := range bldgsRaw {
		if bldgsRaw[i] != "" && addrsRaw[i] != "" && zipsRaw[i] != "" {
			bldgs = append(bldgs, bldgsRaw[i])
			addrs = append(addrs, addrsRaw[i])
			zips = append(zips, zipsRaw[i])
			notes = append(notes, notesRaw[i])
		}
	}

	refDates, refBldgs, _, _, _ := getBldgs()
	regdates := make([]int, 0)
	for _, val := range bldgs {
		tempIndx := findStrInSlice(val, refBldgs)
		if tempIndx != -1 { //if bldg name exist, use existing date
			regdates = append(regdates, refDates[tempIndx])
		} else { //if new bldg name, use current date
			regdates = append(regdates, getCurrentDate())
		}
	}
	intZips := make([]int, 0)
	for _, val := range zips {
		intZip, err := strconv.Atoi(val)
		checkError(err, "controller-updateBldgs")
		intZips = append(intZips, intZip)
	}

	//sort buildings based on bldg name
	sBldgs := bldgs
	sort.Sort(sort.StringSlice(sBldgs))
	sDates := make([]int, len(regdates))
	sAddrs := make([]string, len(addrs))
	sIntZips := make([]int, len(intZips))
	sNotes := make([]string, len(notes))
	for i, _ := range bldgs {
		indx := sort.SearchStrings(sBldgs, bldgs[i])
		sDates[indx] = regdates[i]
		sAddrs[indx] = addrs[i]
		sIntZips[indx] = intZips[i]
		sNotes[indx] = notes[i]
	}

	_ = bldgtableReplace(sDates, sBldgs, sAddrs, sIntZips, sNotes)
	return true
}

func getCustomers() (ids []int, nicknames []string, phones []string, bldgs []string, rooms []string, notes []string) {
	ids, _, nicknames, phonesInt64, bldgs, rooms, notes := customertableGet()
	phonesRaw := make([]string, len(phonesInt64))
	for i, val := range phonesInt64 {
		phonesRaw[i] = strconv.FormatInt(val, 10)
		phones = append(phones, "("+phonesRaw[i][0:3]+") "+phonesRaw[i][3:6]+"-"+phonesRaw[i][6:])
	}
	return
}

func updateCustomers(nicknames []string, phones []string, bldgs []string, rooms []string, notes []string) bool {
	refIds, refDates, refNicknames, refPhones, refBldg, refRooms, refNotes := customertableGet()

	//check any nickname in ref but not in input, remove this row as user changed nickname to ""
	for i, val := range refNicknames {
		if findStrInSlice(val, nicknames) == -1 {
			customertableDelete(refIds[i])
		}
	}

	for i, _ := range phones {
		phones[i] = strings.Replace(phones[i], "(", "", -1) //remove "()- " from []phones
		phones[i] = strings.Replace(phones[i], ")", "", -1)
		phones[i] = strings.Replace(phones[i], "-", "", -1)
		phones[i] = strings.Replace(phones[i], " ", "", -1)
		if len(phones[i]) == 11 { //remove the leading "1" if exists
			phones[i] = phones[i][1:]
		}
		if len(phones[i]) != 10 { //remove invalid phone numbers with length != 10
			phones[i] = ""
		}
		_, err := strconv.ParseInt(phones[i], 10, 64) //remove phone numbers contains non-numbers
		if err != nil {
			phones[i] = ""
		}
		if phones[i] != "" && string(phones[i][0]) == "-" { //remove negative phone numbers
			phones[i] = ""
		}

		if nicknames[i] != "" && phones[i] != "" && bldgs[i] != "" && rooms[i] != "" { //ignore inputs with any of those fields empty
			tempPhone, _ := strconv.ParseInt(phones[i], 10, 64) //guaranteed no error, checked previously
			tempIndx := -1
			for j, _ := range refNicknames {
				if nicknames[i] == refNicknames[j] {
					tempIndx = j
					break
				}
			}

			if tempIndx != -1 { //existing record
				if !(tempPhone == refPhones[tempIndx] && bldgs[i] == refBldg[tempIndx] &&
					rooms[i] == refRooms[tempIndx] && notes[i] == refNotes[tempIndx]) {
					//if not exactly the same, update the row but keep original ID and registration date
					_ = customertableUpdate(refIds[tempIndx], refDates[tempIndx], refNicknames[tempIndx], tempPhone, bldgs[i], rooms[i], notes[i])
				}
			} else { //new record
				customertableAppend(getCurrentDate(), nicknames[i], tempPhone, bldgs[i], rooms[i], notes[i])
			}

		} else if nicknames[i] != "" { //delete record if exists
			for j, _ := range refNicknames {
				if nicknames[i] == refNicknames[j] {
					customertableDelete(refIds[j])
					break
				}
			}
		}
	}

	//now check if any name is being used in active orders, if yes then add it back in with a note
	_, activeNickname, _ := ordertableGetActive()
	if activeNickname != nil {
		uniqueActiveNickname := uniqueStrings(activeNickname)
		_, _, modNickname, _, _, _, _ := customertableGet()

		modNMap := make(map[string]int)
		for _, val := range modNickname {
			modNMap[val] = 1
		}
		for _, val := range uniqueActiveNickname {
			if modNMap[val] == 0 {
				indx := findStrInSlice(val, refNicknames)
				customertableAppend(refDates[indx], val, refPhones[indx], refBldg[indx], refRooms[indx], "Delete aborted. Not able to delete due to active order(s) are under this name")
			}
		}
	}

	return true
}

func logOrders(nicknames []string, items []string, units []string, amounts []string, notes []string) bool {
	nickname := nicknames[0]
	if nickname != "" {
		orderList := ""
		for i, _ := range items {
			if items[i] != "" && units[i] != "" && amounts[i] != "" {
				orderList = orderList + items[i] + "^" + units[i] + "^" + amounts[i] + "^" + notes[i] + "^^^?"
			}
		}
		if orderList != "" {
			orderList = orderList[:len(orderList)-1] //remove the last delimiter
		}
		_ = ordertableAppend(nickname, getCurrentDate(), orderList)
	}
	return true
}

func getOrders() (ids []int, nicknames []string, bldgs []string, orderlists []string) {
	ids, nicknames, orderlists = ordertableGetActive()
	_, bldgNicknames, _, bldgNldgs, _, _ := getCustomers()
	for _, val := range nicknames {
		indx := findStrInSlice(val, bldgNicknames)
		bldgs = append(bldgs, bldgNldgs[indx]) //link the building to every nickname
	}
	return
}

func updateOrders(nickname string, items []string, units []string, amounts []float64, notes []string, units2 []string, amounts2 []float64) {
	//this will remove all active records with the same nickname, then append the new order
	ids, dates := ordertableGetActivePerNickname(nickname)
	tempCheck := 0
	for i := 0; i < len(dates)-1; i++ {
		if dates[i] != dates[i+1] {
			if dates[i] > dates[i+1] { //use the largest date (oldest)
				tempCheck = dates[i]
			} else {
				tempCheck = dates[i+1]
			}
		}
	}
	var date int
	if tempCheck == 0 { //if all dates for this nickname are the same, use it; else use the oldest date
		date = dates[0]
	} else {
		date = tempCheck
	}

	_ = ordertableDelete(ids)

	refItems, refUnits, refUprices, _ := itemtableGet()
	uMap := make(map[string]string)
	pMap := make(map[string]float64)
	for i, _ := range refItems {
		uMap[refItems[i]] = refUnits[i]
		pMap[refItems[i]] = refUprices[i]
	}

	tempOrderList := make([]string, 0)
	for i, _ := range items {
		var price float64
		price = 0
		if uMap[items[i]] == units2[i] { //only record if available unit is found
			price = pMap[items[i]] * amounts2[i]
		}

		tempOrders := []string{items[i], units[i], strconv.FormatFloat(amounts[i], 'f', 2, 64), notes[i], units2[i], strconv.FormatFloat(amounts2[i], 'f', 2, 64), strconv.FormatFloat(price, 'f', 2, 64)}
		tempOrderList = append(tempOrderList, strings.Join(tempOrders, "^"))
	}

	_ = ordertableAppend(nickname, date, strings.Join(tempOrderList, "?"))
	//item list: item ^ unit ^ amount ^ notes ^ unit2 ^ amount2 ^ price? ...
}

func deleteOrders(nicknames []string) {
	idList := make([]int, 0)
	for _, n := range nicknames {
		ids, _ := ordertableGetActivePerNickname(n)
		idList = append(idList, ids...)
	}

	_ = ordertableDelete(idList)
}

func logPurchases(items []string, units []string, amounts []string, prices []string) bool {
	inItems := make([]string, 0)
	inUnits := make([]string, 0)
	inAmounts := make([]float64, 0)
	inPrice := make([]float64, 0)
	for i, _ := range items {
		if items[i] != "" && units[i] != "" && amounts[i] != "" && prices[i] != "" {
			amount, err := strconv.ParseFloat(amounts[i], 64)
			checkError(err, "controller-logPurchases-1")
			price, err := strconv.ParseFloat(prices[i], 64)
			checkError(err, "controller-logPurchases-2")
			inItems = append(inItems, items[i])
			inUnits = append(inUnits, units[i])
			inAmounts = append(inAmounts, amount)
			inPrice = append(inPrice, price)
		}
	}
	_ = purchasetableReplaceActive(getCurrentDate(), inItems, inUnits, inAmounts, inPrice)
	return true
}

func getPurchases() (ids []int, dates []int, items []string, units []string, amounts []float64, prices []float64) {
	ids, dates, items, units, amounts, prices = purchasetableGetActive()
	return
}

func submitTempReport(purchases string, orders string) bool {
	logTempReport(getCurrentDate(), purchases, orders)
	return true
}

func submitReport() bool {
	date, purchases, orders := extractTempReport()
	if logReport(date, purchases, orders) {
		_ = orderpurchaseUpdateStatus()
	}
	return true
}
