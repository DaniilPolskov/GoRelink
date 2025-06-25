package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

var urlMake = make(map[string]string)

func generateRandomID() string {
	number := math.Pow(10, 10)
	randomSeed := rand.Intn(int(number))
	return strconv.Itoa(randomSeed)
}

func main() {
	var oldURL string
	fmt.Print("Enter URL to shorten: ")
	fmt.Scanln(&oldURL)

	if !strings.HasPrefix(oldURL, "http://") && !strings.HasPrefix(oldURL, "https://") {
		oldURL = "https://" + oldURL
	}

	id := generateRandomID()
	shortURL := "http://localhost:8080/gorelink/" + id

	urlMake[id] = oldURL

	fmt.Println("Old URL:", oldURL)
	fmt.Println("Short URL:   ", shortURL)

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
