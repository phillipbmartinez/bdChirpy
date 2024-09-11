package main

import (
	"net/http"
	"log"
)

func main() {
	// Create a new http.ServeMux
	mux := http.NewServeMux();

	// Use http.FileServer to serve files from the current directory
	fileServer := http.FileServer(http.Dir("."))

	// Add a handler for the root path
	mux.Handle("/", fileServer)

	// Create a new http.Server struct
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	// Start the server
	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}