package metrics 

type Metrics interface {
	IncEnqueued()
	IncLeased()
	IncSucceeded()
	IncFailed()
	IncDead()
	SetInProgress(n int)
	SetQueueDepth(n int)
}