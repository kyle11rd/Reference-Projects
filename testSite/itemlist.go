package main

import (
	//"fmt"
	"html/template"
	"net/http"
)

type ItemInfo struct {
	Item      string
	Unit      string
	UnitPrice float64
	Notes     string
	UnitList  []string
}

func itemlist(w http.ResponseWriter, r *http.Request) {
	authCheck(w, r)
	TempUnitList := getUnits()

	TempItems, TempUnits, TempPrice, TempNotes := getItems()

	t, err := template.New("").Funcs(template.FuncMap{
		"isUnit":       isUnit,
		"isValidPrice": isValidPrice,
	}).Parse(tpl_item)
	checkError(err, "itemlist-itemlist-1")
	err = t.ExecuteTemplate(w, "t_top", nil)
	checkError(err, "itemlist-itemlist-2")

	for i, _ := range TempItems {
		iinfo := ItemInfo{
			Item:      TempItems[i],
			Unit:      TempUnits[i],
			UnitPrice: TempPrice[i],
			Notes:     TempNotes[i],
			UnitList:  TempUnitList,
		}
		t, err = template.New("").Funcs(template.FuncMap{
			"isUnit":       isUnit,
			"isValidPrice": isValidPrice,
		}).Parse(tpl_item)
		checkError(err, "itemlist-itemlist-3")
		err = r.ParseForm()
		checkError(err, "itemlist-itemlist-4")
		err = t.ExecuteTemplate(w, "t_info", iinfo)
		checkError(err, "itemlist-itemlist-5")
	}
	for i := 0; i < 5; i++ { //add empty lines to allow adding new customers
		iinfo := ItemInfo{
			Item:      "",
			Unit:      "",
			UnitPrice: -1,
			Notes:     "",
			UnitList:  TempUnitList,
		}
		t, err = template.New("").Funcs(template.FuncMap{
			"isUnit":       isUnit,
			"isValidPrice": isValidPrice,
		}).Parse(tpl_item)
		checkError(err, "itemlist-itemlist-6")
		err = r.ParseForm()
		checkError(err, "itemlist-itemlist-7")
		err = t.ExecuteTemplate(w, "t_info", iinfo)
		checkError(err, "itemlist-itemlist-8")
	}

	t, err = template.New("").Funcs(template.FuncMap{
		"isUnit":       isUnit,
		"isValidPrice": isValidPrice,
	}).Parse(tpl_item)
	checkError(err, "itemlist-itemlist-9")
	err = t.ExecuteTemplate(w, "t_bot", nil)
	checkError(err, "itemlist-itemlist-10")

	if r.Method == "POST" {
		//fmt.Println(r.Form)
		_ = updateItems(r.Form["Item"], r.Form["Unit"], r.Form["Unit Price"], r.Form["Notes"])
	}
}

const tpl_item = `
{{define "t_top"}}
<html>
<head>
<title></title>
<script src="http://code.jquery.com/jquery-1.9.1.js"></script>
<script>
  $(function () {
    $('form').on('submit', function (e) {
      e.preventDefault();
      $.ajax({
        type: 'post',
        data: $('form').serialize(),
      });
    });
  });
</script>
</head>
<body>

<h2>Serving Items</h2>

<form action="/ItemList" method="post">

  <table>

    <tr>
      <td>*Item</td>
      <td>*Unit</td>
      <td>*Unit Price</td>
      <td>Notes</td>
    </tr>
{{end}}

{{define "t_info"}}
<tr>
	<td><input type="text" name="Item" value="{{.Item}}"></td>

	<td>
	    <select name="Unit">
			{{$unit := .Unit}}
			{{range .UnitList}}
				<option {{if isUnit $unit .}}selected{{end}}>{{.}}</option>
			{{end}}
		</select>
	</td>
	
	<td><input type="number" step="0.01" min="0" name="Unit Price" value="{{if isValidPrice .UnitPrice}}{{.UnitPrice}}{{end}}"></td>
	<td><input type="text" name="Notes" value="{{.Notes}}"></td>
</tr>
{{end}}

{{define "t_bot"}}
  </table>

  <br>
  <span>&nbsp</span>
  <input type="submit" value="Update">
</form>
<br>
<a href="/">Click me to go back to main panel</a>

</body>
</html>
{{end}}
`
