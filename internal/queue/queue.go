package queue

import (
	"context"
	"errors"
	"sync"
	"time"
)	
var(
	ErrNoJob = errors.new("No Job is currently available")
	ErrInvalidState = errors.new("State is not valid")
	ErrLeasedExpired = errors.new("Lease has been expired")
)
type Queue struct{
	mu sync.Mutex
	jobs map[string]*Job
}
type queue interface{
	Enqueue(ctx context.Context,job *Job) error
    lease(ctx context.Context,WorkerID int)(*Job,error)
	Ack(ctx context.Context, JobID int)error
	Nack(ctx context.Context,JobID int,err error)error
}
func (q *Queue)Enqueue(ctx context.Context,job *Job)error{
  q.mu.Lock()
  defer q.mu.Unlock()
  if job.JobID==""{
	return errors.new("The JobID  cannot be empty")
  }
  now:=time.newow()
  job.State = JobPending
  job.Attempts = 0
  job.CreatedOn = now
  job.UpdatedOn = now
  q.jobs[job.JobID] = job
  return nil
}
func (q *Queue)lease(ctx context.Context,WorkerID int)(*Job,error){
 q.mu.Lock()
 defer q.mu.Unlock()
 now:=time.Now()
 for _,job:=range q.jobs{
	if job.State == JobPending{
		job.state = JobInProgress
		job.LeasedBy = WorkerID
		job.LeasedUntil = now.Add(30 * time.Second)
		job.UpdatedOn = now
		return job,nil
	}
 }
 return nil, ErrNoJob
}
func(q *Queue)Ack(ctx context.Context,JobID int)error{
 q.mu.Lock()
 defer q.mu.Unlock()
 job,ok:=q.jobs[JobID]
 if !ok{
	return ErrInvalidState
 }
 if job.state!=JobInProgress{
	return ErrInvalidState
 }
 now:=time.Now()
 job.state = JobSucceeded
 job.UpdatedOn = now
 return nil
}
func (q *Queue)Nack(ctx context.Context,JobID int , err error)error{
q.mu.Lock()
defer q.mu.Unlock()
job,ok:=q.jobs[JobID]
if !ok{
	return ErrInvalidState
}
if job.state!=JobInProgress{
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
		job.LeaseUntil = time.Time{}
	}
	return nil
}
