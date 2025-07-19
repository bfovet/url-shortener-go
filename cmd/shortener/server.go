package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

// isValidURL checks if the provided URL is valid
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Start the API server
func Server() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Route("/shorten", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var req Request
			if err := render.DecodeJSON(r.Body, &req); err != nil {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, ErrorResponse{Error: "Invalid JSON"})
				return
			}

			if req.URL == "" {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, ErrorResponse{Error: "URL is required"})
				return
			}

			if !isValidURL(req.URL) {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, ErrorResponse{Error: "Invalid URL format"})
				return
			}

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

	router.Get("/{shortCode}", func(w http.ResponseWriter, r *http.Request) {
		shortCode := chi.URLParam(r, "shortCode")

		originalURL, exists := urlStore[shortCode]

		if !exists {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{Error: "Short URL not found"})
			return
		}

		// Ensure URL has a scheme
		if !strings.HasPrefix(originalURL, "http://") && !strings.HasPrefix(originalURL, "https://") {
			originalURL = "https://" + originalURL
		}

		http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
	})

	fmt.Printf("Server starting on %s\n", defaultHostAddr)
	http.ListenAndServe(defaultHostAddr, router)
}
