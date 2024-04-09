package handler

import (
	"assignment2/myapp/data"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

		var apendix string

		 if allRegistrations[id].Country=="" {
		 	apendix=allRegistrations[id].ISOcode
		 } else if allRegistrations[id].ISOcode=="" {
			response, err := GetContent(data.PATH_RESTCOUNTRIES_API_NAME+allRegistrations[id].Country)
				if err != nil {
	 			log.Fatal(err)
				}
			
			requestApendix:= data.Apendix{}
			json.Unmarshal(response, &requestApendix)
			
			apendix=requestApendix[0].Cca2
		} else {
			apendix=allRegistrations[id].ISOcode
		}

		// Retrievinig data from REST Countries API, latitude, longitude, capital, area, population
		response, err := GetContent(data.PATH_RESTCOUNTRIES_API+apendix)
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

		if len(AllNotification)>0 {
			for i:=0; i<len(AllNotification); i++ {
				if AllNotification[i].Country==apendix {
					if AllNotification[i].Event=="INVOKE" {
						switch r.Method {
							case http.MethodGet: {
	
								currentWebhook:=data.WebhookInvoke{}
								currentWebhook.Country=apendix
								currentWebhook.Event="INVOKE"
								currentWebhook.Id=AllNotification[i].Id
								currentWebhook.Time=formatedTime
	
								currentWebhookJSON, err:= json.Marshal(currentWebhook)
								if err != nil {
									log.Println("Error during request creation. Error:", err)
									return
									}
	
								req, err := http.NewRequest(http.MethodPost, AllNotification[i].Url, bytes.NewBuffer(currentWebhookJSON))
								if err != nil {
									log.Println("Error during request creation. Error:", err)
									return
									}
	
								// Perform invocation
								client1 := http.Client{}
								res, err := client1.Do(req)
								if err != nil {
									log.Println("Error in HTTP request. Error:", err)
									return
									}
	
								// Read the response
								response, err := io.ReadAll(res.Body)
								if err != nil {
									log.Println("Something is wrong with invocation response. Error:", err)
									return
									}
	
								log.Println("Webhook " + AllNotification[i].Url + " invoked. Received status code " + strconv.Itoa(res.StatusCode) +
								" and body: " + string(response))
							}
						default:
							http.Error(w, "Method "+r.Method+" not supported for ", http.StatusMethodNotAllowed)
						}
					}
				}
			}
		}

		request,err:=json.Marshal(responseDashboard)
		if err != nil {
			http.Error(w, "Error during pretty printing", http.StatusInternalServerError)
			return
			}
		
		fmt.Fprintln(w,string(request))
	
	} 
}