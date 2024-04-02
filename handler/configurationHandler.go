package handler

import (
	"assignment2/myapp/data"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var allRegistrations []data.Registration
var idRegistration = 1;

func ConfigurationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ConfigurationsGet(w, r)
	case http.MethodPost:
		ConfigurationsPost(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+" "+http.MethodPost+"' are supported.", http.StatusNotImplemented)
		return
	}
}

func ConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ConfigurationGet(w, r)
	case http.MethodPut:
		ConfigurationPut(w, r)
	case http.MethodDelete:
		ConfigurationDelete(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+" "+http.MethodDelete+" "+http.MethodPut+"' are supported.", http.StatusNotImplemented)
		return
	}
}

func ConfigurationsPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("allow-control-allow-methods", "POST")

	var registration data.Registration
	var currentRegistration data.CurrentRegistration

	err := json.NewDecoder(r.Body).Decode(&registration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Decoding: " + err.Error())
		return
	}

	registration.Id = strconv.Itoa(idRegistration)
	currentRegistration.Id = strconv.Itoa(idRegistration)
	idRegistration++

	t := time.Now()
	formatedTime := t.Format("2006-01-02 15:04:05")
	currentRegistration.Time = formatedTime
	registration.Lastchange = formatedTime

	allRegistrations = append(allRegistrations, registration)

	w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(currentRegistration)
}

func ConfigurationsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("allow-control-allow-methods", "GET")

	var register=len(allRegistrations)

	if register==0 {
		fmt.Fprintln(w, "\n\tThere are no configuration registered. Register configuration first...")
	} else {
		json.NewEncoder(w).Encode(allRegistrations)
	}
}

func ConfigurationGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("allow-control-allow-methods", "GET")

	var register=len(allRegistrations)

	url := strings.Split(r.URL.Path, "/")
	urlValue := url[len(url)-1]

	value, err := strconv.Atoi(urlValue)
    if err != nil {
        panic(err)
    }

	var id=value-1

	if register==0 {
		fmt.Fprintln(w, "\n\tThere are no configuration registered. Register configuration first...")
	} else if value>register {
		fmt.Fprintln(w, "\n\tConfiguration ID is out of range. Dashboard configuration registered:",register)
	} else {
		json.NewEncoder(w).Encode(allRegistrations[id])
	} 
}

func ConfigurationPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.Header().Add("allow-control-allow-methods", "PUT")

	var register=len(allRegistrations)

	url := strings.Split(r.URL.Path, "/")
	urlValue := url[len(url)-1]
 
	value, err := strconv.Atoi(urlValue)
	if err != nil {
		 panic(err)
	}

	var id=value-1

	if register==0 {
		fmt.Fprintln(w, "\n\tThere are no configuration registered. Register configuration first...")
	} else if value>register {
		fmt.Fprintln(w, "\n\tConfiguration ID is out of range. Dashboard configuration registered:",register)
	} else {
		var updateRegistartion data.Registration

		err = json.NewDecoder(r.Body).Decode(&updateRegistartion)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println("Decoding: " + err.Error())
			return
		}

		updateRegistartion.Id = strconv.Itoa(value)

		t := time.Now()
		formatedTime := t.Format("2006-01-02 15:04:05")
		updateRegistartion.Lastchange = formatedTime

		allRegistrations[id]=updateRegistartion
		fmt.Fprintln(w, "\n\tConfiguration with ID", value,"is updated. Dashboard registry updated...")
		}
}

func ConfigurationDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.Header().Add("allow-control-allow-methods", "DELETE")

	var register=len(allRegistrations)

	url := strings.Split(r.URL.Path, "/")
	urlValue := url[len(url)-1]
 
	value, err := strconv.Atoi(urlValue)
	if err != nil {
		 panic(err)
	}

	var id=value-1

	if register==0 {
		fmt.Fprintln(w, "\n\tThere are no configuration registered. Register configuration first...")
	} else if value>register {
		fmt.Fprintln(w, "\n\tConfiguration ID is out of range. Dashboard configuration registered:",register)
	} else {
		fmt.Fprintln(w, "\n\tConfiguration with ID", id+1,"is removed. Dashboard registry updated...")

		allRegistrations = append(allRegistrations[:id],allRegistrations[id+1:]...)

		for i:=0; i<register-1; i++ {
			allRegistrations[i].Id=strconv.Itoa(i+1)
			}
		}

		idRegistration--
		if register==0 {
			idRegistration=1
		}
}