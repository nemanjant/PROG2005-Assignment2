package handler

import (
	"io"
	"math/rand"
	"net/http"
	"time"
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

func GenerateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    seed := rand.NewSource(time.Now().UnixNano())
    random := rand.New(seed)

    result := make([]byte, length)
    for i := range result {
        result[i] = charset[random.Intn(len(charset))]
    }
    return string(result)
}

