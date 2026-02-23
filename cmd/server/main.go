package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
    "strconv"
	"distributed_job_queue/internal/api"
    "distributed_job_queue/internal/metrics"
	"distributed_job_queue/internal/queue"
)

func main() {
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	q := &queue.MemoryQueue{}
	m := metrics.New()

	handlers := &api.Handlers{
		Queue:   q,
		Metrics: m,
	}

	

	mux := http.NewServeMux()
	api.RegisterRoutes(mux, handlers)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	

	workerCount := 4
	for i := 0; i < workerCount; i++ {
		workerID := "worker-" + strconv.Itoa(i)

		go func(id string) {
			w := queue.NewWorker(id, q)
			w.Run(ctx)
		}(workerID)
	}

	

	go func() {
		log.Println("HTTP server listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	log.Println("shutdown signal received")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("shutdown complete")
}