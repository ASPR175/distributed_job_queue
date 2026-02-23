package queue

import (
	"context"
	"errors"
	"sync"
	"time"
)	
var(
	ErrNoJob = errors.New("No Job is currently available")
	ErrInvalidState = errors.New("State is not valid")
	ErrLeasedExpired = errors.New("Lease has been expired")
)
type MemoryQueue struct {
	mu   sync.Mutex
	jobs map[string]*Job
}
type Queue interface {
	Enqueue(ctx context.Context, job *Job) error
	Lease(ctx context.Context, WorkerID string) (*Job, error)
	Ack(ctx context.Context, jobID string) error
	Nack(ctx context.Context, jobID string, err error) error
}
func (q *MemoryQueue)Enqueue(ctx context.Context,job *Job)error{
  q.mu.Lock()
  defer q.mu.Unlock()
  if job.JobID==""{
	return errors.New("The JobID  cannot be empty")
  }
  now:=time.Now()
  job.State = JobPending
  job.Attempts = 0
  job.CreatedOn = now
  job.UpdatedOn = now
  if q.jobs == nil {
	q.jobs = make(map[string]*Job)
  }
  q.jobs[job.JobID] = job
  return nil
}
func (q *MemoryQueue) Lease(ctx context.Context, WorkerID string) (*Job, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	now := time.Now()
	for _, job := range q.jobs {
		if job.State == JobPending ||
   (job.State == JobInProgress && now.After(job.LeasedUntil)) {
			job.State = JobInProgress
			job.LeasedBy = WorkerID
			job.LeasedUntil = now.Add(30 * time.Second)
			job.UpdatedOn = now
			return job, nil
		}
	}
	return nil, ErrNoJob
}
func(q *MemoryQueue)Ack(ctx context.Context,JobID string)error{
 q.mu.Lock()
 defer q.mu.Unlock()
 job,ok:=q.jobs[JobID]
 if !ok{
	return ErrInvalidState
 }
if job.LeasedBy == ""{
    return ErrInvalidState
}
 now:=time.Now()
 job.State = JobSucceeded
 job.UpdatedOn = now
 return nil
}
func (q *MemoryQueue)Nack(ctx context.Context,JobID string , err error)error{
q.mu.Lock()
defer q.mu.Unlock()
job,ok:=q.jobs[JobID]
if !ok{
	return ErrInvalidState
}
if job.State!=JobInProgress{
	return ErrInvalidState
}
job.Attempts++
job.LastError = err.Error()
job.UpdatedOn = time.Now()
if job.Attempts >= job.MaxAttempts {
		job.State = JobDead
	} else {
		job.State = JobPending
		job.LeasedBy = ""
		job.LeasedUntil = time.Time{}
	}
	return nil
}
