package main

import (
	"html/template"
	"net/http"
)

type UnitInfo struct {
	Units []string
}

func unitslist(w http.ResponseWriter, r *http.Request) {
	authCheck(w, r)
	ulist := UnitInfo{
		Units: getUnits(),
	}

	for i := 0; i < 5; i++ { //leave emply slots for new units
		ulist.Units = append(ulist.Units, "")
	}

	t, err := template.New("").Parse(tpl_unit)
	checkError(err, "unitslist-unitslist-1")
	err = r.ParseForm()
	checkError(err, "unitslist-unitslist-2")
	err = t.Execute(w, ulist)
	checkError(err, "unitslist-unitslist-3")

	if r.Method == "POST" {
		_ = updateUnits(r.Form["Unit"])
	}
}

const tpl_unit = `
<html>
<head>
<title>Update Units</title>
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
<h2>Units</h2>
<form method="post">

  <table>

    <tr>
      <td>Units</td>
    </tr>

{{range .Units}}
    <tr>
      <td><input type="text" name="Unit" value="{{.}}"></td>
    </tr>
{{end}}
  </table>

  <br>
  <span>&nbsp</span>
  <input type="submit" value="Save">
</form>
<br>
<a href="/">Click me to go back to main panel</a>
</body>
</html>
`
