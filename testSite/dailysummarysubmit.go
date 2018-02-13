package main

import (
	//"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Report struct {
	Pitems   []string
	Punits   []string
	Pamounts []string
	Pprices  []string

	Oitems   []string
	Ounits   []string
	Oamounts []string
	Oprices  []string
}

func dailysummarysubmit(w http.ResponseWriter, r *http.Request) {
	authCheck(w, r)
	if r.Method == "GET" {
		url := r.URL.String()
		sIndx := strings.Index(url, "?")
		if sIndx == -1 { //if no selection, print something
			t, _ := template.New("").Parse(tpl_noSelection)
			_ = t.Execute(w, "")
		} else {
			url = url[sIndx+1:]
			url = strings.Replace(url, "%20", " ", -1)
			selections := strings.Split(url, "+")
			_, _, bldgs, orders := getOrders()
			uniqueBldgs := uniqueStrings(bldgs)
			if len(selections) != len(uniqueBldgs) { //in case direct copy&paste the url without all bldgs
				t, _ := template.New("").Parse(tpl_wrong)
				_ = t.Execute(w, "")
			} else {
				_, _, pItems, pUnits, pAmounts, pPrices := getPurchases()
				oItemsUnits := make([]string, 0) //use unit2 as final unit
				oAmounts := make([]float64, 0)   //use amount2 as final amount
				oPrices := make([]float64, 0)
				for _, val := range orders {
					stuff := strings.Split(val, "?")
					for _, val2 := range stuff {
						stuff2 := strings.Split(val2, "^")
						oItemsUnits = append(oItemsUnits, stuff2[0]+"@"+stuff2[4])
						tempAmount, _ := strconv.ParseFloat(stuff2[5], 64)
						tempPrice, _ := strconv.ParseFloat(stuff2[6], 64)
						oAmounts = append(oAmounts, tempAmount)
						oPrices = append(oPrices, tempPrice)
					}
				}
				//pItems, pUnits, pAmounts, pPrice, oItemsUnits, oAmounts, oPrices
				oMap := make(map[string][]int)
				oIUunique := make([]string, 0)
				for i, val := range oItemsUnits {
					if oMap[val] == nil {
						oMap[val] = []int{i}
						oIUunique = append(oIUunique, val)
					} else {
						oMap[val] = append(oMap[val], i)
					}
				}

				oAmod := make([]float64, 0)
				oPmod := make([]float64, 0)

				for _, val := range oIUunique {
					var tempAmount, tempPrice float64
					tempAmount = 0
					tempPrice = 0
					for _, val2 := range oMap[val] {
						tempAmount += oAmounts[val2]
						tempPrice += oPrices[val2]
					}
					oAmod = append(oAmod, tempAmount)
					oPmod = append(oPmod, tempPrice)
				}
				//pItems, pUnits, pAmounts, pPrice, oIUunique, oAmod, oPmod
				reportList := Report{}
				purchases := make([]string, 0)
				orders := make([]string, 0)
				for i, _ := range pItems {
					reportList.Pitems = append(reportList.Pitems, pItems[i])
					reportList.Punits = append(reportList.Punits, pUnits[i])
					reportList.Pamounts = append(reportList.Pamounts, strconv.FormatFloat(pAmounts[i], 'f', 2, 64))
					reportList.Pprices = append(reportList.Pprices, strconv.FormatFloat(pPrices[i], 'f', 2, 64))
					tempList := []string{reportList.Pitems[i], reportList.Punits[i], reportList.Pamounts[i], reportList.Pprices[i]}
					purchases = append(purchases, strings.Join(tempList, ","))
				}
				for i, val := range oIUunique {
					indx := strings.Index(val, "@")
					reportList.Oitems = append(reportList.Oitems, val[0:indx])
					reportList.Ounits = append(reportList.Ounits, val[indx+1:])
					reportList.Oamounts = append(reportList.Oamounts, strconv.FormatFloat(oAmod[i], 'f', 2, 64))
					reportList.Oprices = append(reportList.Oprices, strconv.FormatFloat(oPmod[i], 'f', 2, 64))
					tempList := []string{reportList.Oitems[i], reportList.Ounits[i], reportList.Oamounts[i], reportList.Oprices[i]}
					orders = append(orders, strings.Join(tempList, ","))
				}

				t, err := template.New("").Parse(tpl_report)
				checkError(err, "dailysummarysubmit-dailysummarysubmit-1")
				err = t.Execute(w, reportList)
				checkError(err, "dailysummarysubmit-dailysummarysubmit-2")
				_ = submitTempReport(strings.Join(purchases, ";"), strings.Join(orders, ";"))
			}
		}
	}
	if r.Method == "POST" {
		_ = submitReport()
		t, _ := template.New("").Parse(tpl_success)
		_ = t.Execute(w, "")
	}
}

const tpl_wrong = `
<html>
<body>
<h1>Don't hack me!!!</h1>
</body>
</html>
`

const tpl_success = `
<html>
<body>
<br>
<h2>Report submitted successfully :)</h2>
<br><br>
<a href="/">Click me to go back to main panel</a>
</body>
</html>
`

const tpl_report = `
<html>
<head>
<style>
table, th, td {
    border: 1px solid black;
    border-collapse: collapse;
    padding: 5px;
}
</style>
<h2>Submit Report</h2>
</head>
<body>

<table>
  <tr>
    <th>Item Purchased</th>
	<th>Unit</th>
    <th>Amount</th>		
    <th>Spent</th>
  </tr>
  {{range $i, $e := .Pitems}}
    <tr>
      <td>{{.}}</td>
	  <td>{{index $.Punits $i}}</td>
      <td>{{index $.Pamounts $i}}</td>		
      <td>{{index $.Pprices $i}}</td>
    </tr>
  {{end}}
</table>
<br><br>
<table>
  <tr>
    <th>Item Ordered</th>
    <th>Unit Assigned</th>		
    <th>Amount Delivered</th>
	<th>Revenue</th>
  </tr>
  {{range $i, $e := .Oitems}}
    <tr>
      <td>{{.}}</td>
	  <td>{{index $.Ounits $i}}</td>
      <td>{{index $.Oamounts $i}}</td>
      <td>{{index $.Oprices $i}}</td>
    </tr>
  {{end}}
</table>
<br>
<form method="post">
<input type="submit" value="Submit Report">
</form>
</body>
</html>
`
