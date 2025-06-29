package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var (
	defaultHostAddr = ":8080"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
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
		})
	})

	http.ListenAndServe(defaultHostAddr, router)
}
