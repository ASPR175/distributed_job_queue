package metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type Metrics interface {
	IncEnqueued()
	IncLeased()
	IncSucceeded()
	IncFailed()
	IncDead()

	SetInProgress(n int)
	SetQueueDepth(n int)

	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type InMemoryMetrics struct {
	enqueued    int64
	leased      int64
	succeeded   int64
	failed      int64
	dead        int64
	inProgress  int64
	queueDepth  int64
}

func New() Metrics {
	return &InMemoryMetrics{}
}

func (m *InMemoryMetrics) IncEnqueued()  { atomic.AddInt64(&m.enqueued, 1) }
func (m *InMemoryMetrics) IncLeased()    { atomic.AddInt64(&m.leased, 1) }
func (m *InMemoryMetrics) IncSucceeded() { atomic.AddInt64(&m.succeeded, 1) }
func (m *InMemoryMetrics) IncFailed()    { atomic.AddInt64(&m.failed, 1) }
func (m *InMemoryMetrics) IncDead()      { atomic.AddInt64(&m.dead, 1) }

func (m *InMemoryMetrics) SetInProgress(n int) {
	atomic.StoreInt64(&m.inProgress, int64(n))
}

func (m *InMemoryMetrics) SetQueueDepth(n int) {
	atomic.StoreInt64(&m.queueDepth, int64(n))
}
func (m *InMemoryMetrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,
		`jobs_enqueued %d
jobs_leased %d
jobs_succeeded %d
jobs_failed %d
jobs_dead %d
jobs_in_progress %d
queue_depth %d
`,
		atomic.LoadInt64(&m.enqueued),
		atomic.LoadInt64(&m.leased),
		atomic.LoadInt64(&m.succeeded),
		atomic.LoadInt64(&m.failed),
		atomic.LoadInt64(&m.dead),
		atomic.LoadInt64(&m.inProgress),
		atomic.LoadInt64(&m.queueDepth),
	)
}