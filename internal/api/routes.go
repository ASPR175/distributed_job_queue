package api

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, h *Handlers) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/jobs", h.EnqueueJob)
	mux.HandleFunc("/metrics", h.Metrics)
}