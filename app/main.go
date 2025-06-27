package main

import (
	"log"
	"net/http"

	"GoRelink/handlers"
	"GoRelink/storage"
)

func main() {
	store := storage.NewMemoryStore()

	http.HandleFunc("/api/shorten", handlers.ShortenHandler(store))
	http.HandleFunc("/api/shorten/batch", handlers.BatchShortenHandler(store))
	http.HandleFunc("/gorelink/", handlers.RedirectHandler(store))

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
