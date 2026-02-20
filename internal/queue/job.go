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
func(j JobState) String()string{
	switch j{
	case JobPending: return "Pending"
	case JobInProgress: return "In_Progress"
	case JobSucceeded: return "Success"
	case JobFailed: return "Failed"
	case JobDead : return "Dead"
	default: return "Unknown"
	}
}
type Job struct{
	JobID string
	Payload []byte
	State JobState
	Attempts int
	MaxAttempts int
	LeasedBy string
	LeasedUntil time.Time
	LastError string
	CreatedOn time.Time
	UpdatedOn time.Time
}