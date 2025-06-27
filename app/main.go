package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"GoRelink/shortener"
	"GoRelink/storage"
	"GoRelink/types"
)

func withCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func generateUniqueID(store *storage.MemoryStore) string {
	for {
		id := shortener.GenerateID()
		if _, exists := store.Get(id); !exists {
			return id
		}
	}
}

func normalizeURL(raw string) string {
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		return "https://" + raw
	}
	return raw
}

func main() {
	store := storage.NewMemoryStore()

	http.HandleFunc("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		withCORS(w)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req types.ShortenRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.URL == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		req.URL = normalizeURL(req.URL)

		id := generateUniqueID(store)

		store.Save(id, req.URL)

		shortURL := fmt.Sprintf("http://localhost:8080/gorelink/%s", id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"shortURL": shortURL,
		})
	})

	http.HandleFunc("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		withCORS(w)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req types.BatchShortenRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" || req.Count <= 0 || req.Count > 10 {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		req.URL = normalizeURL(req.URL)

		shortURLs := make([]string, 0, req.Count)
		done := make(chan string, req.Count)

		for i := 0; i < req.Count; i++ {
			go func() {
				id := generateUniqueID(store)
				store.Save(id, req.URL)
				short := fmt.Sprintf("http://localhost:8080/gorelink/%s", id)
				done <- short
			}()
		}

		for i := 0; i < req.Count; i++ {
			shortURLs = append(shortURLs, <-done)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]string{
			"shortURLs": shortURLs,
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
