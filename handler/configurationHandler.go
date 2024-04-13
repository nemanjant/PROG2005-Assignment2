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

// Storage of all registrations during service run, id auto increment/decrement
var allRegistrations []data.Registration
var idRegistration = 1;

// Switch between differnet methods for given handler, new configuration and all configurations
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

// Switch between differnet methods for given handler, specific configuration
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

// Function to create new dashboard register POST Method
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

	// Webhook notification on REGISTER event
	if len(AllNotification)>0 {
		for i:=0; i<len(AllNotification); i++ {
			if AllNotification[i].Country==registration.ISOcode {
				if AllNotification[i].Event=="REGISTER" {
					switch r.Method {
						case http.MethodPost: {

							currentWebhook:=data.WebhookInvoke{}
							currentWebhook.Country=registration.ISOcode
							currentWebhook.Event="REGISTER"
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
							client := http.Client{}
							res, err := client.Do(req)
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

	w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(currentRegistration)
}

// Function to retrieve all registered configurations GET Method
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

// Function to retrieve specific registered configurations GET Method
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

// Function to update specific registered configurations PUT Method
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

		// Webhook notification on CHANGE event
		if len(AllNotification)>0 {
			for i:=0; i<len(AllNotification); i++ {
				if AllNotification[i].Country==allRegistrations[id].ISOcode {
					if AllNotification[i].Event=="CHANGE" {
						switch r.Method {
							case http.MethodPut: {

								currentWebhook:=data.WebhookInvoke{}
								currentWebhook.Country=AllNotification[i].Country
								currentWebhook.Event="CHANGE"
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
								client := http.Client{}
								res, err := client.Do(req)
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
		fmt.Fprintln(w, "\n\tConfiguration with ID", value,"is updated. Dashboard registry updated...")
		}
}

// Function to delete specific registered configurations DELETE Method
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

		t := time.Now()
		formatedTime := t.Format("2006-01-02 15:04:05")

		// Webhook notification on DELETE event
		if len(AllNotification)>0 {
			for i:=0; i<len(AllNotification); i++ {
				if AllNotification[i].Country==allRegistrations[id].ISOcode {
					if AllNotification[i].Event=="DELETE" {
						switch r.Method {
							case http.MethodDelete: {

								currentWebhook:=data.WebhookInvoke{}
								currentWebhook.Country=AllNotification[i].Country
								currentWebhook.Event="DELETE"
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
								client := http.Client{}
								res, err := client.Do(req)
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