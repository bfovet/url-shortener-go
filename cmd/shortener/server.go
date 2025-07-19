package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var (
	defaultHostAddr = ":8080"
	urlStore        = make(map[string]string) // shortCode -> originalURL
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
}

type Response struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// generateShortCode creates a random 6-character code
func generateShortCode() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:6]
}

// Start the API server
func Server() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Route("/shorten", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var req Request
			err := render.DecodeJSON(r.Body, &req)
			if err != nil {
				fmt.Println("failed to decode request body")
				return
			}
			fmt.Println(req)

			// Generate unique short code
			var shortCode string
			for {
				shortCode = generateShortCode()
				if _, exists := urlStore[shortCode]; !exists {
					urlStore[shortCode] = req.URL
					break
				}
			}

			baseURL := fmt.Sprintf("http://localhost%s", defaultHostAddr)
			shortURL := fmt.Sprintf("%s/%s", baseURL, shortCode)

			response := Response{
				ShortCode: shortCode,
				ShortURL:  shortURL,
			}

			render.Status(r, http.StatusCreated)
			render.JSON(w, r, response)
		})
	})

	fmt.Printf("Server starting on %s\n", defaultHostAddr)
	http.ListenAndServe(defaultHostAddr, router)
}
