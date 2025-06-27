package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"GoRelink/shortener"
	"GoRelink/storage"
)

func main() {
	store := storage.NewMemoryStore()

	var oldURL string
	fmt.Print("Enter URL to shorten: ")
	fmt.Scanln(&oldURL)

	if !strings.HasPrefix(oldURL, "http://") && !strings.HasPrefix(oldURL, "https://") {
		oldURL = "https://" + oldURL
	}

	var id string
	for {
		id = shortener.GenerateID()
		if _, exists := store.Get(id); !exists {
			break
		}
	}

	store.Save(id, oldURL)

	shortURL := "http://localhost:8080/gorelink/" + id

	fmt.Println("Old URL:", oldURL)
	fmt.Println("Short URL:", shortURL)

	http.HandleFunc("/gorelink/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/gorelink/")
		oldURL, found := store.Get(id)
		if !found {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, oldURL, http.StatusFound)
	})

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
