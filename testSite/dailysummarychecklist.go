package main

import (
	//"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func dailysummarychecklist(w http.ResponseWriter, r *http.Request) {
	authCheck(w, r)
	url := r.URL.String()
	sIndx := strings.Index(url, "?")
	if sIndx == -1 { //if no selection, print something
		t, _ := template.New("").Parse(tpl_noSelection)
		_ = t.Execute(w, "")
	} else {
		url = url[sIndx+1:]
		url = strings.Replace(url, "%20", " ", -1)
		selections := strings.Split(url, "+")

		_, refNicknames, refBldgs, refOrders := getOrders()
		//make a map to improve speed
		refMap := make(map[string][]int)
		for i, n := range refNicknames {
			if refMap[n] == nil {
				refMap[n] = []int{i}
			} else {
				refMap[n] = append(refMap[n], i)
			}
		}

		_, cName, cPhone, _, cRoom, _ := getCustomers()
		cusMap := make(map[string]int)
		for i, n := range cName {
			cusMap[n] = i
		}

		t, err := template.New("").Parse(tpl_checklist)
		checkError(err, "dailysummarychecklist-dailysummarychecklist-1")
		err = t.ExecuteTemplate(w, "t_top", "")
		checkError(err, "dailysummarychecklist-dailysummarychecklist-2")

		//first sort bldg list & nickname from orders
		for _, bldg := range selections { //rotate over selections (bldgs)
			//sort nicknames per selection (bldg)
			tempList := make([]string, 0) //take nicknames per selection (bldg)
			for i, refBldg := range refBldgs {
				if refBldg == bldg {
					tempList = append(tempList, refNicknames[i])
				}
			}
			bldgName := BldgName{
				BldgName: bldg,
			}
			t, err = template.New("").Parse(tpl_checklist)
			checkError(err, "dailysummarychecklist-dailysummarychecklist-3")
			err = t.ExecuteTemplate(w, "t_bldgTop", bldgName)
			checkError(err, "dailysummarychecklist-dailysummarychecklist-4")
			tempList = uniqueStrings(tempList)
			sort.Strings(tempList) //sort unique nicknames
			for _, tempNickname := range tempList {
				refIndx := refMap[tempNickname]
				for _, i := range refIndx {
					//write out orders sorted by nickname per bldg
					tempOrderList := strings.Split(refOrders[i], "?")
					var total float64
					total = 0
					for _, val := range tempOrderList {
						tempItemList := strings.Split(val, "^")
						tempNum, _ := strconv.ParseFloat(tempItemList[6], 64)
						total += tempNum
					}
					for indx, val := range tempOrderList {
						tempItemList := strings.Split(val, "^")
						tempFloat, _ := strconv.ParseFloat(tempItemList[2], 64)
						tempFloat2, _ := strconv.ParseFloat(tempItemList[5], 64)
						tempFloat3, _ := strconv.ParseFloat(tempItemList[6], 64)

						iInfo := OrderStuc{
							Nickname: tempNickname,
							Phone:    cPhone[cusMap[tempNickname]],
							Room:     cRoom[cusMap[tempNickname]],
							Span:     1,
							Item:     tempItemList[0],
							Unit:     tempItemList[1],
							Amount:   tempFloat,
							Note:     tempItemList[3],
							Unit2:    tempItemList[4],
							Amount2:  tempFloat2,
							Price:    tempFloat3,
							Total:    strconv.FormatFloat(total, 'f', 2, 64),
							IsFirst:  false,
						}

						if indx == 0 {
							iInfo.IsFirst = true
							iInfo.Span = len(tempOrderList)
						}

						t, err := template.New("").Parse(tpl_checklist)
						checkError(err, "dailysummarychecklist-dailysummarychecklist-7")
						err = t.ExecuteTemplate(w, "t_loop", iInfo)
						checkError(err, "dailysummarychecklist-dailysummarychecklist-8")
					}
				}
			}
			t, err = template.New("").Parse(tpl_checklist)
			checkError(err, "dailysummarychecklist-dailysummarychecklist-11")
			err = t.ExecuteTemplate(w, "t_bldgBot", "")
			checkError(err, "dailysummarychecklist-dailysummarychecklist-12")
		}
		t, err = template.New("").Parse(tpl_checklist)
		checkError(err, "dailysummarychecklist-dailysummarychecklist-13")
		err = t.ExecuteTemplate(w, "t_bot", "")
		checkError(err, "dailysummarychecklist-dailysummarychecklist-14")
	}
}

//width=670px is save to fit letter sized papers
const tpl_checklist = `
{{define "t_top"}}
<html>
<head>

<style>
table, th, td {
    border: 1px solid black;
    border-collapse: collapse;
}
th, td {
    padding: 5px;
    text-align: left;    
}
.A {
	width: 150px;
	max-width: 150px;
}
.B {
	width: 100px;
	max-width: 100px;
}
.C {
	width: 50px;
	max-width: 50px;
}
.D {
	width: 50px;
	max-width: 50px;
}
.E {
	width: 100px;
	max-width: 100px;
	word-wrap: normal;
}
.F {
	width: 50px;
	max-width: 50px;
}
.G {
	width: 55px;
	max-width: 55px;
}
.H {
	width: 60px;
	max-width: 60px;
}
.I {
	width: 60px;
	max-width: 60px;
	text-align: center;
}
.notbold{
    font-weight:normal
}​
</style>
</head>
<body>

<form method="post">
{{end}}

{{define "t_bldgTop"}}
<table>
<h3>{{.BldgName}}</h3>
  <tr>
    <th class="A">Nickname</th>
    <td class="B">Item</td>
    <td class="C">Unit</td>
    <td class="D">Amount</td>
	<td class="E">Notes</td>
    <td class="F">Unit*</td>
    <td class="G">Amount*</td>
	<td class="H">Price [$]</td>
	<th class="I">Total [$]</th>
  </tr>
{{end}}

{{define "t_loop"}}
<tr>
    {{if .IsFirst}}
	  <th rowspan="{{.Span}}" class="A">{{.Nickname}}
	  &#160;&#160;&#160;&#160;&#160;&#160;&#160;&#160;<span class="notbold">{{.Phone}}</span><br>
	  &#160;&#160;&#160;&#160;&#160;&#160;&#160;&#160;<span class="notbold">Room: {{.Room}}</span></th>
	{{end}}
    <td class="B">{{.Item}}</td>
    <td class="C">{{.Unit}}</td>
    <td class="D">{{.Amount}}</td>
	<td class="E">{{.Note}}</td>
    <td class="F">{{.Unit2}}</td>
    <td class="G">{{.Amount2}}</td>
	<td class="H">{{.Price}}</td>
	{{if .IsFirst}}
	  <th rowspan="{{.Span}}" class="I">{{.Total}}</th>
	{{end}}
</tr>
{{end}}

{{define "t_bldgBot"}}
</table>
<br>
{{end}}

{{define "t_bot"}}
</form>
</body>
</html>
{{end}}
`
