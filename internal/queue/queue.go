package queue

import "context"
var(
	ErrNoJob = errors.new("No Job is currently available")
	ErrInvalidState = errors.new("State is not valid")
	ErrLeasedExpired = errors.new("Lease has been expired")
)
type queue interface{
	Enqueue(ctx context.Context,job *Job) error
    lease(ctx context.Context,WorkerID int)(*Job,error)
	Ack(ctx context.Context, JobID int)error
	Nack(ctx context.Context,JobID int,err error)error
}
