package main

import (
    "fmt"
    "log"
    "net/http"
)

type apiConfig struct {
    fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits++
        next.ServeHTTP(w, r)
    })
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    _, err := w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
    if err != nil {
        log.Printf("Error writing response: %v", err)
    }
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
    cfg.fileserverHits = 0
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    _, err := w.Write([]byte("Hits reset"))
    if err != nil {
        log.Printf("Error writing response: %v", err)
    }
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    _, err := w.Write([]byte("OK"))
    if err != nil {
        log.Printf("Error writing response: %v", err)
    }
}

func main() {
    const filepathRoot = "."
    const port = "8080"

    apiCfg := &apiConfig{}

    mux := http.NewServeMux()

    // Serve files from the root directory under /app/ path
    fileServer := http.FileServer(http.Dir(filepathRoot))
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", fileServer)))

    // Register the readiness handler
    mux.HandleFunc("/healthz", handlerReadiness)
    
    // Register the metrics handler
    mux.HandleFunc("/metrics", apiCfg.handleMetrics)
    
    // Register the reset handler
    mux.HandleFunc("/reset", apiCfg.handleReset)

    srv := &http.Server{
        Addr:    ":" + port,
        Handler: mux,
    }

    log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
    log.Fatal(srv.ListenAndServe())
}