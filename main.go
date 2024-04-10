package main

import (
	"assignment2/myapp/data"
	"assignment2/myapp/handler"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	http.HandleFunc(data.PATH_REGISTRATIONS, handler.ConfigurationsHandler)
	http.HandleFunc(data.PATH_REGISTRATION_ID, handler.ConfigurationHandler)
	http.HandleFunc(data.PATH_DASHBOARD_ID, handler.DashboardHandler)
	http.HandleFunc(data.PATH_NOTIFICATIONS, handler.NotificationsHandler)
	http.HandleFunc(data.PATH_NOTIFICATIONS_ID, handler.NotificationHandler)
	http.HandleFunc(data.PATH_STATUS, handler.StatusHandler)

	log.Println("Starting server on port " + port + "...")
	log.Fatal(http.ListenAndServe(":" + port,nil))
}
