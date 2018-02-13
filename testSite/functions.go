package main

import (
	"strconv"
	"strings"
	"time"
)

func findStrInSlice(str string, slice []string) int {
	for i, val := range slice {
		if str == val {
			return i
		}
	}
	return -1
}

func getCurrentDate() int {
	currentTime := time.Now().Local()
	currentDate := currentTime.Format("20060102")
	intDate, err := strconv.Atoi(currentDate)
	checkError(err, "functions-getCurrentDate")
	return intDate
}

func getCurrentDateStr() string {
	currentTime := time.Now().Local()
	currentDate := currentTime.Format("20060102")
	wDate := []string{currentDate[4:6], currentDate[6:8], currentDate[2:4]}
	return strings.Join(wDate, "/")
}

func isValidPrice(price float64) bool {
	return price > 0
}

func isUnit(unit string, unitOption string) bool {
	return unit == unitOption
}

func isMatching(a string, b string) bool {
	return a == b
}

func uniqueStrings(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}
	for v := range elements {
		if encountered[elements[v]] == false {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

//for future use
/*func uniqueS(elements []string) (result []string, occurance []int) {
	encountered := map[string]int{}
	for v := range elements {
		if encountered[elements[v]] == 0 {
			encountered[elements[v]] = 1
			result = append(result, elements[v])
		} else {
			encountered[elements[v]] += 1
		}
	}
	for _, n := range result {
		occurance = append(occurance, encountered[n])
	}
	return
}*/

func checkboxParser(inList []string) (status []int) {
	//checkbox will carry [0 0 1 0] as string for [off on off]
	cnt := 0
	for cnt < len(inList)-1 {
		if inList[cnt] == "0" && inList[cnt+1] == "1" {
			status = append(status, 1)
			cnt += 2
		} else {
			status = append(status, 0)
			cnt += 1
		}
	}
	if inList[len(inList)-1] == "0" {
		status = append(status, 0)
	}
	return
}
