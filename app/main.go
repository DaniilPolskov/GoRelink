package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"GoRelink/shortener"
	"GoRelink/storage"
)

func main() {
	store := storage.NewMemoryStore()

	http.HandleFunc("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			URL string `json:"url"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.URL == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
			req.URL = "https://" + req.URL
		}

		var id string
		for {
			id = shortener.GenerateID()
			if _, exists := store.Get(id); !exists {
				break
			}
		}

		store.Save(id, req.URL)

		shortURL := fmt.Sprintf("http://localhost:8080/gorelink/%s", id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"shortURL": shortURL,
		})
	})

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
