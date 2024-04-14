package handler

import (
	"assignment2/myapp/data"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Function to retrieve and close given URL
func GetContent(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

// Function to create random string value
func IdGenerator(idlength int) string {
    const char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    start := rand.NewSource(time.Now().UnixNano())
    new := rand.New(start)

    id := make([]byte, idlength)
    for i := range id {
        id[i] = char[new.Intn(len(char))]
    }
    return string(id)
}

// Function to create random string value, time independent
func RandString(n int) string {
    const alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    var bytes = make([]byte, n)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = alphanum[b % byte(len(alphanum))]
    }
    return string(bytes)
}

// Function to retrieve stored notificications from Firestore
func ReadFirestore(n data.Notification, m []data.Notification) (mn []data.Notification){
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

	iter := client.Collection(data.COLLECTION).Documents(ctx)
	defer iter.Stop() 
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
		return
		}

		if err := doc.DataTo(&n); err != nil {
			log.Println(err)
			return
		}
		m = append(m, n)
	}

    return m
}
