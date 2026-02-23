package api

import (
	"encoding/json"
	"net/http"
	"time"

	"distributed_job_queue/internal/queue"
	"distributed_job_queue/internal/metrics"
)
type Handlers struct {
	Queue   queue.Queue
	Metrics metrics.Metrics
}
type enqueueRequest struct {
	ID          string `json:"id"`
	Payload     []byte `json:"payload"`
	MaxAttempts int    `json:"max_attempts"`
}
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
func (h *Handlers) EnqueueJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req enqueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	job := &queue.Job{
		JobID:       req.ID,
		Payload:     req.Payload,
		MaxAttempts: req.MaxAttempts,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.Queue.Enqueue(ctx, job); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.Metrics.IncEnqueued()

	w.WriteHeader(http.StatusAccepted)
}
func (h *Handlers) Metrics(w http.ResponseWriter, r *http.Request) {
	h.Metrics.ServeHTTP(w, r)
}