package main

import (
	"html/template"
	"net/http"
)

type CustomerInfo struct {
	Nickname string
	Phone    string
	Bldg     string
	Room     string
	Notes    string
	BldgList []string
}

func customerlist(w http.ResponseWriter, r *http.Request) {
	authCheck(w, r)
	_, TempBldgList, _, _, _ := getBldgs()
	_, TempNicknames, TempPhone, TempBldg, TempRoom, TempNotes := getCustomers()
	t, err := template.New("").Funcs(template.FuncMap{
		"isBldg": isBldg,
	}).Parse(tpl_customer)
	checkError(err, "customerlist-customerlist-1")
	err = t.ExecuteTemplate(w, "t_top", nil)
	checkError(err, "customerlist-customerlist-2")

	for i, _ := range TempNicknames {
		cinfo := CustomerInfo{
			Nickname: TempNicknames[i],
			Phone:    TempPhone[i],
			Bldg:     TempBldg[i],
			Room:     TempRoom[i],
			Notes:    TempNotes[i],
			BldgList: TempBldgList,
		}
		t, err = template.New("").Funcs(template.FuncMap{
			"isBldg": isBldg,
		}).Parse(tpl_customer)
		checkError(err, "customerlist-customerlist-3")
		err = r.ParseForm()
		checkError(err, "customerlist-customerlist-4")
		err = t.ExecuteTemplate(w, "t_info", cinfo)
		//err = t.ExecuteTemplate(w, "t_info", cinfo)
		checkError(err, "customerlist-customerlist-5")
	}
	for i := 0; i < 10; i++ { //add empty lines to allow adding new customers
		cinfo := CustomerInfo{
			Nickname: "",
			Phone:    "",
			Bldg:     "",
			Room:     "",
			Notes:    "",
			BldgList: TempBldgList,
		}
		t, err = template.New("").Funcs(template.FuncMap{
			"isBldg": isBldg,
		}).Parse(tpl_customer)
		checkError(err, "customerlist-customerlist-6")
		err = r.ParseForm()
		checkError(err, "customerlist-customerlist-7")
		err = t.ExecuteTemplate(w, "t_info", cinfo)
		//err = t.ExecuteTemplate(w, "t_info", cinfo)
		checkError(err, "customerlist-customerlist-8")
	}

	t, err = template.New("").Funcs(template.FuncMap{
		"isBldg": isBldg,
	}).Parse(tpl_customer)
	checkError(err, "customerlist-customerlist-9")
	err = t.ExecuteTemplate(w, "t_bot", nil)
	checkError(err, "customerlist-customerlist-10")

	if r.Method == "POST" {
		//fmt.Println(r.Form)
		_ = updateCustomers(r.Form["Nickname"], r.Form["Phone"], r.Form["Building"], r.Form["Room"], r.Form["Notes"])
	}
}

func isBldg(bldg string, bldgOption string) bool {
	return bldg == bldgOption
}

const tpl_customer = `
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

<h2>Customer List</h2>

<form action="/CustomerList" method="post">

  <table>

    <tr>
      <td>*Nickname</td>
      <td>*Phone Number</td>
      <td>*Building</td>
      <td>*Room</td>
      <td>Notes</td>
    </tr>
{{end}}

{{define "t_info"}}
<tr>
	<td><input type="text" name="Nickname" value="{{.Nickname}}"></td>
	<td><input type="text" name="Phone" value="{{.Phone}}"></td>
	
	<td>
	    <select name="Building">
			{{$bldg := .Bldg}}
			{{range .BldgList}}
				<option {{if isBldg $bldg .}}selected{{end}}>{{.}}</option>
			{{end}}
		</select>
	</td>
	
	<td><input type="text" name="Room" value="{{.Room}}"></td>
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
