package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/shirwahersi/go-s3-browser/internal/s3"
)

type Handler struct {
	s3Client *s3.Client
}

func NewHandler(s3Client *s3.Client) *Handler {
	return &Handler{
		s3Client: s3Client,
	}
}

// ListHandler handles GET /api/list?prefix=<path>
func (h *Handler) ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	prefix := r.URL.Query().Get("prefix")

	result, err := h.s3Client.ListObjects(r.Context(), prefix)
	if err != nil {
		log.Printf("Error listing objects: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to list objects"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// DownloadHandler handles GET /api/download?key=<file-key>
func (h *Handler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing key parameter"})
		return
	}

	url, err := h.s3Client.GetPresignedURL(r.Context(), key, time.Hour)
	if err != nil {
		log.Printf("Error generating download URL: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate download URL"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}
