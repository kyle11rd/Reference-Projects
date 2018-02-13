package main

import (
	//"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Receipts struct {
	NicknameNDate string
	Items         []string
	AmountsNUnits []string
	Prices        []string
	Total         string
}

const LINES_PER_PAGE int = 50 //number of rows to fit 1 letter sized paper
const COL_PER_PAGE int = 3    //number of columns to fit 1 letter sized paper
const HEADER_LINES int = 4    //number of rows for rows other than item/amount/price

func dailysummaryreceipt(w http.ResponseWriter, r *http.Request) {
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

		t, _ := template.New("").Parse(tpl_receipt)
		_ = t.ExecuteTemplate(w, "t_top", "")
		t, _ = template.New("").Parse(tpl_receipt)
		_ = t.ExecuteTemplate(w, "t_newPageTop", "")
		t, _ = template.New("").Parse(tpl_receipt)
		_ = t.ExecuteTemplate(w, "t_newColumnTop", "")

		colCnt := 1
		rowCnt := 1
		dateStr := getCurrentDateStr()
		for _, bldg := range selections { //rotate over selections (bldgs)
			//sort nicknames per selection (bldg)
			tempList := make([]string, 0) //take nicknames per selection (bldg)
			for i, refBldg := range refBldgs {
				if refBldg == bldg {
					tempList = append(tempList, refNicknames[i])
				}
			}

			for _, n := range tempList {
				refIndx := refMap[n]
				for _, i := range refIndx {
					tempOrderList := strings.Split(refOrders[i], "?")
					tempNickname := n
					if len(n) > 20 {
						tempNickname = tempNickname[:20]
					}
					inData := Receipts{
						NicknameNDate: tempNickname + " - " + dateStr,
						Items:         make([]string, 0),
						AmountsNUnits: make([]string, 0),
						Prices:        make([]string, 0),
						Total:         "",
					}

					var tempTotal float64
					tempTotal = 0
					tempRowCnt := HEADER_LINES
					for _, val := range tempOrderList {
						tempItemList := strings.Split(val, "^")
						inData.Items = append(inData.Items, tempItemList[0])
						inData.AmountsNUnits = append(inData.AmountsNUnits, tempItemList[5]+" "+tempItemList[4])
						inData.Prices = append(inData.Prices, tempItemList[6])
						tempFloat, err := strconv.ParseFloat(tempItemList[6], 64)
						checkError(err, "dailysummaryreceipt-dailysummaryreceipt-1")
						tempTotal = tempTotal + tempFloat
						tempRowCnt += 1
					}
					inData.Total = strconv.FormatFloat(tempTotal, 'f', 2, 64)

					if tempRowCnt+rowCnt > LINES_PER_PAGE { //switch column and/or page if needed

						//close current column
						t, _ := template.New("").Parse(tpl_receipt)
						_ = t.ExecuteTemplate(w, "t_newColumnBot", "")
						colCnt += 1

						//if colCnt > COL_PER_PAGE, move to a new page
						if colCnt > COL_PER_PAGE {
							colCnt = 1
							t, _ = template.New("").Parse(tpl_receipt)
							_ = t.ExecuteTemplate(w, "t_newPageBot", "")
							t, _ = template.New("").Parse(tpl_receipt)
							_ = t.ExecuteTemplate(w, "t_newLine", "")
							t, _ = template.New("").Parse(tpl_receipt)
							_ = t.ExecuteTemplate(w, "t_newPageTop", "")
						}

						//at last start a new column
						t, _ = template.New("").Parse(tpl_receipt)
						_ = t.ExecuteTemplate(w, "t_newColumnTop", "")
						rowCnt = 1
					}
					t, err := template.New("").Parse(tpl_receipt)
					checkError(err, "dailysummaryreceipt-dailysummaryreceipt-2")
					err = t.ExecuteTemplate(w, "t_newRecord", inData)
					checkError(err, "dailysummaryreceipt-dailysummaryreceipt-3")
					rowCnt += tempRowCnt

				}
			}
		}

		t, _ = template.New("").Parse(tpl_receipt)
		_ = t.ExecuteTemplate(w, "t_newColumnBot", "")
		t, _ = template.New("").Parse(tpl_receipt)
		_ = t.ExecuteTemplate(w, "t_newPageBot", "")
		t, _ = template.New("").Parse(tpl_receipt)
		_ = t.ExecuteTemplate(w, "t_bot", "")
	}

}

/*
Nickname - mmddyy
Item  Amount&Unit Price
...   ...         ...
-----------------------
Total: xxx
*/

const tpl_receipt = `
{{define "t_top"}}
<html>
<head>
<style>
table{
	border-collapse: collapse;
}
p{
	text-align: center;
	margin: 0px;
	max-width: 225px;
}
.item{
	display: block;
	min-width: 100px;
	max-width: 100px;
	text-align: left;
}
.anu{
	display: block;
	min-width: 70px;
	max-width: 70px;
	text-align: left;
}
.price{
	display: block;
	min-width: 50px;
	max-width: 50px;
	text-align: right;
}
.total{
	display: block;
	min-width: 225px;
	max-width: 225px;
	border-top: 1px solid black;
	text-align: right;
}
.totalprice{
	display: inline-block;
	min-width: 60px;
	max-width: 60px;
	text-align: right;
}

.namedate {
	display: inline-block;
}

.mainTd {
	padding-right: 20px;
	height: 910px;
}

</style>
</head>
<body>
{{end}}

{{define "t_newPageTop"}}
<table>
<tr>
{{end}}

{{define "t_newColumnTop"}}
	<td class="mainTd">
{{end}}

{{define "t_newRecord"}}
		<p><b><i>~&#160&#160&#160GoodOBag&#160&#160&#160~</i></b></p>
		<span class="namedate">{{.NicknameNDate}}</span>
		<table>
			{{range $i, $e := .Items}}
			<tr>
				<td><span class="item">{{.}}</span></td>
				<td><span class="anu">{{index $.AmountsNUnits $i}}</span></td>
				<td class="price"><span>{{index $.Prices $i}}</span></td>
			</tr>
			{{end}}
		</table>
		<span class="total">Total [$] :<span class="totalprice">{{.Total}}</span></span>
		<br>
{{end}}

{{define "t_newColumnBot"}}
	</td>
{{end}}

{{define "t_newPageBot"}}
</tr>
</table>
{{end}}

{{define "t_bot"}}
</body>
</html>
{{end}}
`
