package main

import (
	"assignment2/myapp/data"
	"assignment2/myapp/handler" // Firestore-specific support
	"context"                   // State handling across API boundaries; part of native GoLang API
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go" // Generic firebase support
	"google.golang.org/api/option"
)

const collection = "messages"

func main() {

	ctx:= context.Background()

	opt := option.WithCredentialsFile("./assignment2-8c8dd-firebase-adminsdk-1q43z-b1f562cd40.json")
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
	result,err:=ref.Set(ctx,map[string]interface{}{
		"url": "https://localhost:8080/client/", 
   		"country": "NO",                         
   		"event": "INVOKE",                        
	})
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




	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	http.HandleFunc(data.PATH_REGISTRATIONS, handler.ConfigurationsHandler)
	http.HandleFunc(data.PATH_REGISTRATION_ID, handler.ConfigurationHandler)
	http.HandleFunc(data.PATH_DASHBOARD_ID, handler.DashboardHandler)

	log.Println("Starting server on port " + port + "...")
	log.Fatal(http.ListenAndServe(":" + port,nil))
}