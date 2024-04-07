package handler

import (
	"assignment2/myapp/data"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Storage of all registered notifications during service run
var AllNotification []data.Notification

const collection = "messages"

// Switch between differnet methods for given handler, new configuration and all configurations
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		NotificationsPost(w, r)
	case http.MethodGet:
	 	NotificationsGet(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+" "+http.MethodPost+"' are supported.", http.StatusNotImplemented)
		return
	}
}

// Switch between differnet methods for given handler, specific configuration
func NotificationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
	 	NotificationGet(w, r)
	case http.MethodDelete:
	 	NotificationDelete(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+" "+http.MethodDelete+"' are supported.", http.StatusNotImplemented)
		return
	}
}

// Function to create new notification register POST Method
func NotificationsPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("allow-control-allow-methods", "POST")

	var notification data.Notification
	var CurrentNotification data.CurrentNotification

	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Decoding: " + err.Error())
		return
	}

	notification.Id=GenerateRandomString(13)

	CurrentNotification.Id=notification.Id

	AllNotification = append(AllNotification, notification)

	// Sending registered notification to Firebase/Firestorm databse
	ctx:= context.Background()

	opt := option.WithCredentialsFile("./credentials/assignment2credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v", err)
  		return 
		}

	client, err:= app.Firestore(ctx)
	if err != nil {
		log.Println(err)
		return
		}

	ref:=client.Collection(collection).NewDoc()
	result,err:=ref.Set(ctx,notification)
	if err != nil {
		log.Println(err)
		return
	 	}
	
	log.Printf("Result is [%v]", result)
		
	defer func() {
		errClose := client.Close()
		if errClose != nil {
			log.Fatal("Closing of the Firebase client failed. Error:", errClose)
		}
	}()

	w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(CurrentNotification)
}

// Function to retrieve all registered notifications GET Method
func NotificationsGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("allow-control-allow-methods", "GET")

	var register=len(AllNotification)

	if register==0 {
		fmt.Fprintln(w, "\n\tThere are no notifications registered. Register notification first...")
	} else {
		json.NewEncoder(w).Encode(AllNotification)
	}
}

// Function to retrieve specific registered notification GET Method
func NotificationGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("allow-control-allow-methods", "GET")

	var register=len(AllNotification)

	url := strings.Split(r.URL.Path, "/")
	urlValue := url[len(url)-1]
 
	if register==0 {
		fmt.Fprintln(w, "\n\tThere are no notification registered. Register notification first...")
	} else {
		for i:=0; i<register; i++ {
			if AllNotification[i].Id==urlValue {
				json.NewEncoder(w).Encode(AllNotification[i])
			}
		}
	}
}

// Function to delete specific registered notification DELETE Method
func NotificationDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.Header().Add("allow-control-allow-methods", "DELETE")

	var register=len(AllNotification)

	url := strings.Split(r.URL.Path, "/")
	urlValue := url[len(url)-1]

	if register==0 {
		fmt.Fprintln(w, "\n\tThere are no notification registered. Register notification first...")
	} else {
		for i:=0; i<register; i++ {
			if AllNotification[i].Id==urlValue {
				AllNotification = append(AllNotification[:i],AllNotification[i+1:]...)
				register--
				fmt.Fprintln(w, "\n\tNotification with ID '",urlValue,"' is removed. Notification registry updated...")
			}
		} 
	}
}