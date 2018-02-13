package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Subdomains map[string]http.Handler

const DOMAIN_BODY = "localhost:8080"

func (subdomains Subdomains) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domainParts := strings.Split(r.Host, ".")
	mux := domainParts[0]
	if mux == "www" { //remove the effect of www
		mux = DOMAIN_BODY
	}
	if subdomains[mux] != nil { //subdomain
		subdomains[mux].ServeHTTP(w, r)
	} else {
		if mux == DOMAIN_BODY { //domain
			subdomains[""].ServeHTTP(w, r)
		} else {
			http.Error(w, "Sorry, page is not found", 404)
		}
	}
}

func main() {
	//domain
	mainMux := http.NewServeMux()
	mainMux.HandleFunc("/", indxPage)

	//subdomain "manage"
	manageMux := http.NewServeMux()
	manageMux.HandleFunc("/login", login)
	manageMux.HandleFunc("/", manage)
	manageMux.HandleFunc("/OrderInfo", orderinfo)                         //Log Orders
	manageMux.HandleFunc("/PurchaseInfo", purchaseinfo)                   //Log Purchases
	manageMux.HandleFunc("/UnitsList", unitslist)                         //lb, ea, etc.
	manageMux.HandleFunc("/ItemList", itemlist)                           //List of items to serve
	manageMux.HandleFunc("/CustomerList", customerlist)                   //List of customers
	manageMux.HandleFunc("/BuildingList", buildinglist)                   //Serving buildings
	manageMux.HandleFunc("/DailySummary", dailysummary)                   //Daily Summary selection page
	manageMux.HandleFunc("/DailySummaryRecords", dailysummaryrecords)     //Daily Summary summary page
	manageMux.HandleFunc("/DailySummaryPrint", dailysummaryprint)         //Daily Summary printout page
	manageMux.HandleFunc("/DailySummaryChecklist", dailysummarychecklist) //Daily Summary checklist for print
	manageMux.HandleFunc("/DailySummaryReceipt", dailysummaryreceipt)     //Daily Summary receipts for print
	manageMux.HandleFunc("/DailySummarySubmit", dailysummarysubmit)       //Daily Summary submittion

	//default settings
	subdomains := make(Subdomains)
	subdomains[""] = mainMux
	subdomains["manage"] = manageMux

	err := http.ListenAndServe(":8080", subdomains)
	checkError(err, "main-main-subdomain")
}

func indxPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to GoodOBag. For any inquiry, please contact goodobag@gmail.com")
}

func checkError(err error, loc string) {
	//format of loc:  scriptName-functionName-anyOtherComment
	//e.g. loc for error in current function = main-checkError
	if err != nil {
		log.Fatal(loc, ":", err)
	}
}
