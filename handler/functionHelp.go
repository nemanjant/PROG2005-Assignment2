package handler

import (
	"io"
	"net/http"
)

func ToFloat32(in int) float32 {
 	return float32(in)
}

// Function to retrieve and close given URL
func GetContent(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}