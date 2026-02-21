package retry


import (
	"math/rand"
	"time"
)

type Backoff interface {
	Next(attempt int) time.Duration
}

type ExponentialBackoff struct {
	BaseDelay time.Duration
	MaxDelay  time.Duration
	Jitter    bool
}

func (b ExponentialBackoff) Next(attempt int) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}

	
	delay := b.BaseDelay << (attempt - 1)

	if delay > b.MaxDelay {
		delay = b.MaxDelay
	}

	if b.Jitter {
		
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		delay = delay/2 + jitter
	}

	return delay
}