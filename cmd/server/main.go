package main

import (
	"log"
	"net/http"
	"time"

	"document-extractor/internal/ai"
	"document-extractor/internal/config"
	"document-extractor/internal/handler"
	"document-extractor/internal/parser"
	"document-extractor/internal/service"
)

func main() {
	cfg := config.Load()

	client := ai.NewClient(cfg.OllamaURL, cfg.OllamaModel, cfg.OllamaTimeout)
	extractor := ai.NewDocumentExtractor(client)
	parsers := parser.NewRegistry()
	documentService := service.NewDocumentService(parsers, extractor)
	h := handler.New(documentService, cfg.MaxFileSize)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /api/v1/extract", h.Extract)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           loggingMiddleware(mux),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       cfg.HTTPTimeout,
		WriteTimeout:      cfg.HTTPTimeout,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("document extractor listening on %s", cfg.HTTPAddr)
	log.Printf("ollama: %s, model: %s", cfg.OllamaURL, cfg.OllamaModel)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
