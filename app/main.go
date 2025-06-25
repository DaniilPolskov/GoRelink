package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

var urlMake = make(map[string]string)

const idLength = 6
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomID() string {
	b := make([]byte, idLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	var oldURL string
	fmt.Print("Enter URL to shorten: ")
	fmt.Scanln(&oldURL)

	if !strings.HasPrefix(oldURL, "http://") && !strings.HasPrefix(oldURL, "https://") {
		oldURL = "https://" + oldURL
	}

	var id string
	for {
		id = generateRandomID()
		if _, exists := urlMake[id]; !exists {
			break
		}
	}

	shortURL := "http://localhost:8080/gorelink/" + id

	urlMake[id] = oldURL

	fmt.Println("Old URL:", oldURL)
	fmt.Println("Short URL:", shortURL)

	http.HandleFunc("/gorelink/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/gorelink/"):]
		oldURL, found := urlMake[id]
		if !found {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, oldURL, http.StatusFound)
	})

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
