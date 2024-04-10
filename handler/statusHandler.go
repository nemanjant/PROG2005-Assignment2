package handler

import (
	"assignment2/myapp/data"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var start = time.Now()

//Acquiring status of given endpoints;
func GetStatusCode(url string) int {
	resp, _ := http.Get(url)
	return resp.StatusCode
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		StatusHandlerGet(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+"' are supported.", http.StatusNotImplemented)
		return
	}
}

func StatusHandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.Header().Add("method", "GET")

	status := data.Status{
		//Retrieving data and checkin if all endpoints er active, using "no" (Norway) value;
		//If no value is represented API returns code 204 and 404;
		// If all is correct returns 200 OK;
		ContriesApi: GetStatusCode(data.PATH_RESTCOUNTRIES_API + "no"),
		MeteoApi: GetStatusCode(data.PATH_OPENMETEO_NO_STATUS),
		CurrencyApi: GetStatusCode(data.PATH_CURRENCY_API + "nok"),
		//Current version of API;
		Version:    data.VERSION,
		Webhooks:   len(AllNotification),
		UpTime:     time.Now().Sub(start).Seconds(),
	}

	responseStruct,err:=json.Marshal(status)
	if err != nil {
		http.Error(w, "Error during pretty printing", http.StatusInternalServerError)
		return
		}

	fmt.Println(string(responseStruct))
	fmt.Fprintln(w, string(responseStruct))
}