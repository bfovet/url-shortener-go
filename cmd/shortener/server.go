package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ShortenUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /shorten")
	io.WriteString(w, "shortening url")
}

// Start the API server
func Server() {
	http.HandleFunc("/shorten", ShortenUrl)

	err := http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
