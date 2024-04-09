package handler

import (
	"io"
	"math/rand"
	"net/http"
	"time"
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
