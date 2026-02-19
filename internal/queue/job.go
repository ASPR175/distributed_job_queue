package queue

import "time"
type JobState int

const(
	JobPending JobState = iota
	JobInProgress 
	JobSucceeded
	JobFailed
	JobDead
)
func(j JobState) string()string{
	switch j{
	case JobPending: return "Pending"
	case JobInProgress: return "In_Progress"
	case JobSucceeded: return "Success"
	case JobFailed: return "Failed"
	case JobDead : return "Dead"
	}
}
type Job struct{
	JobID int
	Payload []byte
	State JobState
	Attempts int
	MaxAttempts int
	LeasedBy string
	LeasedUntil time.time
	LastError string
	CreatedOn time.time
	UpdatedOn time.time
}