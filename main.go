package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/shirwahersi/go-s3-browser/internal/config"
	"github.com/shirwahersi/go-s3-browser/internal/handlers"
	"github.com/shirwahersi/go-s3-browser/internal/s3"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize S3 client
	s3Client, err := s3.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	// Parse index.html as template
	indexHTML, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		log.Fatalf("Failed to read index.html: %v", err)
	}
	indexTemplate, err := template.New("index").Parse(string(indexHTML))
	if err != nil {
		log.Fatalf("Failed to parse index.html template: %v", err)
	}

	// Create handlers
	handler := handlers.NewHandler(s3Client)

	// Set up routes
	http.HandleFunc("/api/list", handler.ListHandler)
	http.HandleFunc("/api/download", handler.DownloadHandler)

	// Serve static files from embedded filesystem (except index.html)
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to create static filesystem: %v", err)
	}
	fileServer := http.FileServer(http.FS(staticFS))

	// Custom handler to serve index.html as template, other files as static
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve index.html as template for root path or if file doesn't exist
		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if err := indexTemplate.Execute(w, cfg.Site); err != nil {
				log.Printf("Failed to execute template: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		// For other paths, check if file exists and serve static content
		path := strings.TrimPrefix(r.URL.Path, "/")
		if _, err := staticFiles.ReadFile("static/" + path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// File not found - serve index.html for SPA routing
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := indexTemplate.Execute(w, cfg.Site); err != nil {
			log.Printf("Failed to execute template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
