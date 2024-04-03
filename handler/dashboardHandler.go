package handler

import (
	"assignment2/myapp/data"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Switch between differnet methods for given handler
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		DashboardGet(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+"' are supported.", http.StatusNotImplemented)
		return
	}
}

// Function to retrieve registered dashboard configurations GET Method
func DashboardGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.Header().Add("method", "GET")

	register := len(allRegistrations)

	url := strings.Split(r.URL.Path, "/")
	urlValue := url[len(url)-1]

	value, err := strconv.Atoi(urlValue)
    if err != nil {
        panic(err)
    	}

	var id = value-1

	if register==0 {
		fmt.Fprintln(w, "\n\tThere are no dashboard configuration registered. Register dashboard first...")
	} else if value>register {
		fmt.Fprintln(w, "\n\tDashboard ID is out of range. Dashboard configuration registered:",register)
	} else {
		responseDashboard:= data.DashboardResponse{}

	// Retrievinig data from REST Countries API, latitude, longitude, capital, area, population
	response, err := GetContent(data.PATH_RESTCOUNTRIES_API+allRegistrations[id].ISOcode)
		if err != nil {
	 	log.Fatal(err)
		}

	requestCountry:= data.CountryRequest{}
	json.Unmarshal(response, &requestCountry)

	var capital = requestCountry[0].Capital[0]
	var latitude = requestCountry[0].Latlng[0]
	var longitude = requestCountry[0].Latlng[1]
	var currency string
	
	for key:= range requestCountry[0].Currencies {
		currency = key
	}

	responseDashboard.ISOcode=allRegistrations[id].ISOcode
	responseDashboard.Country=allRegistrations[id].Country

	// If stattmets for optional imput in configuration registration, bool values
	if (allRegistrations[id].Features.Precipitation) {

		latstring:= fmt.Sprintf("%f", latitude)
		lonstring:= fmt.Sprintf("%f", longitude)
		
		responseMeteo:=data.MeteoRequest{}

		// Retrievinig data from Open-Meteo APIs, precipitation
		response, err := GetContent("https://api.open-meteo.com/v1/forecast?latitude="+latstring+"&longitude="+lonstring+"&daily=precipitation_sum&forecast_days=1")
			if err != nil {
	 		log.Fatal(err)
			}

		json.Unmarshal(response, &responseMeteo)

		responseDashboard.Features.Precipitation=float32(responseMeteo.Daily.PrecipitationSum[0])
	} 

	if (allRegistrations[id].Features.Temperature) {

		latstring:= fmt.Sprintf("%f", latitude)
		lonstring:= fmt.Sprintf("%f", longitude)
		
		responseMeteo:=data.MeteoRequest{}

		// Retrievinig data from Open-Meteo APIs, temperature
		response, err := GetContent("https://api.open-meteo.com/v1/forecast?latitude="+latstring+"&longitude="+lonstring+"&daily=temperature_2m_max,temperature_2m_min&forecast_days=1")
			if err != nil {
	 		log.Fatal(err)
			}

		json.Unmarshal(response, &responseMeteo)

		temperature2MMin:=float32(responseMeteo.Daily.Temperature2MMin[0])
		temperature2MMax:=float32(responseMeteo.Daily.Temperature2MMax[0])

		averageTemperature2M:=(temperature2MMin+temperature2MMax) / 2

		responseDashboard.Features.Temperature=averageTemperature2M
	} 

	if (allRegistrations[id].Features.Capital) {
		responseDashboard.Features.Capital=capital
	}

	if (allRegistrations[id].Features.Coordinates) {
		responseDashboard.Features.Coordinates.Latitude=latitude
		responseDashboard.Features.Coordinates.Longitude=longitude
	} 

	if (allRegistrations[id].Features.Population) {
		responseDashboard.Features.Population=requestCountry[0].Population
	} 

	if (allRegistrations[id].Features.Area) {
		responseDashboard.Features.Area=requestCountry[0].Area
	}

	// Retrievinig data from Currency API
	if len(allRegistrations[id].Features.TargetCurrencies)>0 {
		response, err := GetContent(data.PATH_CURRENCY_API+currency)
			if err != nil {
	 		log.Fatal(err)
			}
		
		responseCurrency:=data.CurrencyRequest{}
		
		json.Unmarshal(response, &responseCurrency)

		currentCurrencies:= make(map[string]float32)

		// Checking currencies for given
		for _, i:=range allRegistrations[id].Features.TargetCurrencies {
			for j, k:= range responseCurrency.Rates {
				if j==i {
					currentCurrencies[j]=k
				}
			}
		}

		responseDashboard.Features.TargetCurrencies=currentCurrencies
	}

	t := time.Now()
	formatedTime := t.Format("2006-01-02 15:04:05")
	responseDashboard.Lastchange = formatedTime

	request,err:=json.Marshal(responseDashboard)
	if err != nil {
		http.Error(w, "Error during pretty printing", http.StatusInternalServerError)
		return
		}
	
	fmt.Fprintln(w,string(request))
	} 
}