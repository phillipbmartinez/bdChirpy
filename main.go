package main

import (
	"net/http"
	"log"
)

// readinessHandler handles requests to /healthz
func readinessHandler(w http.ResponseWriter, r *http.Request) {
    // Set the Content-Type header
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

    // Set the status code to 200 OK
    w.WriteHeader(http.StatusOK)

    // Write the body text
    _, err := w.Write([]byte("OK"))
    if err != nil {
        log.Printf("Error writing response: %v", err)
    }
}

func main() {
	// Create a new http.ServeMux
	mux := http.NewServeMux();

	// Register the readiness endpoint
    mux.HandleFunc("/healthz", readinessHandler)

	// Create a file server for the root directory
    fileServer := http.FileServer(http.Dir("."))

	// Handle requests to /app/ with the file server
    mux.Handle("/app/", http.StripPrefix("/app/", fileServer))

	// Create a file server for the 'assets' directory
	assetsDir := http.Dir("assets")
	assetsFileServer := http.FileServer(assetsDir)

	// Handle requests to /app/assets/ with the assets file server
    mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", assetsFileServer))

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