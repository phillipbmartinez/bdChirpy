package main

import (
	"net/http"
	"log"
)

func main() {
	// Create a new http.ServeMux
	mux := http.NewServeMux();

	/// Handle the root path to serve index.html
    mux.Handle("/", http.FileServer(http.Dir(".")))

	// Create a file server for the 'assets' directory
	assetsDir := http.Dir("assets")
	fileServer := http.FileServer(assetsDir)

	// Handle requests to /assets/ with the file server
    mux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

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