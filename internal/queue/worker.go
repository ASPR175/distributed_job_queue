package queue


import (
	"context"
	"log"
	"time"
)

type JobHandler func(ctx context.Context, job *Job) error

type Worker struct {
	WorkerID      string
	Queue   Queue
	Handler JobHandler

	PollInterval time.Duration
}

func (w *Worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %s shutting down", w.WorkerID)
			return
		default:
		}

		job, err := w.Queue.Lease(ctx, w.WorkerID)
		if err != nil {
			if err == ErrNoJob {
				time.Sleep(w.PollInterval)
				continue
			}

			log.Printf("worker %s lease error: %v", w.WorkerID, err)
			time.Sleep(w.PollInterval)
			continue
		}

		err = w.Handler(ctx, job)
		if err != nil {
			_ = w.Queue.Nack(ctx, job.JobID, err)
			continue
		}

		_ = w.Queue.Ack(ctx, job.JobID)
	}
}