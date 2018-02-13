package main

import (
	"html/template"
	"net/http"
	"strings"
)

type ReadyForSubmit struct {
	Msg     string
	IsValid bool
}

func dailysummaryprint(w http.ResponseWriter, r *http.Request) {
	authCheck(w, r)
	url := r.URL.String()
	defaultUrl := url
	sIndx := strings.Index(url, "?")
	if sIndx == -1 { //if no selection, print something
		t, _ := template.New("").Parse(tpl_noSelection)
		_ = t.Execute(w, "")
	} else {
		//first check if selected all bldgs and all amount are not 0
		//if yes then show the submit report button

		//first check if all bldgs are selected
		_, _, bldgs, _ := getOrders()
		uniqueBldgs := uniqueStrings(bldgs)
		url = url[sIndx+1:]
		url = strings.Replace(url, "%20", " ", -1)
		selections := strings.Split(url, "+")
		isReady := ReadyForSubmit{}
		isReady.Msg = ""
		if len(selections) != len(uniqueBldgs) {
			isReady.Msg = "To submit daily report, please select and review orders from all buildings"
		}

		//now check if all amount are not 0
		if isReady.Msg == "" {
			_, nicknames, _, orderList := getOrders()
			for i, val := range orderList {
				tempList1 := strings.Split(val, "?")
				for _, stuff := range tempList1 {
					tempList2 := strings.Split(stuff, "^")
					if tempList2[6] == "0.00" { //by default precision of price == 0.01
						isReady.Msg = "Please validate UnitAssigned and AmountAssigned for nickname " + nicknames[i]
						break
					}
				}
				if isReady.Msg != "" {
					break
				}
			}
		}
		if isReady.Msg == "" {
			isReady.IsValid = true
		} else {
			isReady.IsValid = false
		}

		if r.Method == "GET" {
			t, err := template.New("").Parse(tpl_print)
			checkError(err, "dailysummaryprint-dailysummaryprint-1")
			err = t.Execute(w, isReady)
			checkError(err, "dailysummaryprint-dailysummaryprint-2")
		}
		if r.Method == "POST" {
			err := r.ParseForm()
			checkError(err, "dailysummaryprint-dailysummaryprint-3")
			if r.Form["print"][0] == "Print Checklist" {
				http.Redirect(w, r, "/DailySummaryChecklist"+defaultUrl[sIndx:], http.StatusSeeOther)
			} else if r.Form["print"][0] == "Print Receipt" {
				http.Redirect(w, r, "/DailySummaryReceipt"+defaultUrl[sIndx:], http.StatusSeeOther)
			} else if r.Form["print"][0] == "Submit Daily Report" {
				http.Redirect(w, r, "/DailySummarySubmit"+defaultUrl[sIndx:], http.StatusSeeOther)
			}

		}
	}
}

const tpl_print = `
<html>
<head>
<style>
input[type=submit] {
  width: 200px; 
  height: 50px;
  font-size: 20px;
}
</style>
<h3>Please note the data may take a few seconds to process</h3>
</head>
<body>
<form method="post">
<input type="submit" name="print" value="Print Checklist">
<br><br>
<input type="submit" name="print" value="Print Receipt">
<br><br>
{{if .IsValid}}
<input type="submit" name="print" value="Submit Daily Report">
{{else}}
<h3>{{.Msg}}</h3>
{{end}}
</form>
<br>
<a href="/">Click me to go back to main panel</a>
</body>
</html>
`
