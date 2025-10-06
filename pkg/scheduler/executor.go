package scheduler

import (
	"context"
	"log"
	"sync"
)

// Tracks concurrent state
type Executor struct {
	ctx    context.Context
	cancel context.CancelCauseFunc

	mu          sync.Mutex
	done        map[string]bool
	failed      bool
	executed    []string
	status      map[string]*JobStatus
	subscribers []chan []JobStatus
}

func (ex *Executor) markDone(id string) {
	ex.mu.Lock()
	defer ex.mu.Unlock()
	ex.done[id] = true
	ex.executed = append(ex.executed, id)
}

func (ex *Executor) markFailed(id string) {
	ex.mu.Lock()
	defer ex.mu.Unlock()
	ex.failed = true
}

func (ex *Executor) isReady(d *SchedulerImpl, id string) bool {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	parents, err := d.dag.GetParents(id)
	if err != nil {
		return false
	}
	for pid := range parents {
		if !ex.done[pid] {
			return false
		}
	}
	return true
}

func (ex *Executor) revertAll(d *SchedulerImpl) {
	ex.mu.Lock()
	executed := make([]string, len(ex.executed))
	copy(executed, ex.executed)
	ex.mu.Unlock()

	for i := len(executed) - 1; i >= 0; i-- {
		id := executed[i]
		jobVal, err := d.dag.GetVertex(id)
		if err != nil {
			continue
		}
		job, ok := jobVal.(Job)
		if !ok {
			continue
		}
		log.Println("Revert:", id)
		if err := job.Revert(ex.ctx); err != nil {
			log.Println("Revert failed for", id, ":", err)
			continue
		}
		ex.setState(id, JobReverted, nil)
	}
	log.Println("Revert complete.")
}

func (ex *Executor) setState(id string, state JobState, err error) {
	ex.mu.Lock()
	s := ex.status[id]
	if s == nil {
		s = &JobStatus{ID: id}
		ex.status[id] = s
	}
	s.State = state
	s.Error = err

	// Copy subscribers and statuses while holding the lock
	statuses := make([]JobStatus, 0, len(ex.status))
	for _, s := range ex.status {
		statuses = append(statuses, *s)
	}
	subs := append([]chan []JobStatus(nil), ex.subscribers...)

	ex.mu.Unlock()

	// Send updates outside the lock
	for _, sub := range subs {
		select {
		case sub <- statuses:
		default:
			// drop if subscriber isn't reading fast enough
		}
	}
}

func (ex *Executor) subscribe() <-chan []JobStatus {
	ch := make(chan []JobStatus, 1)
	ex.mu.Lock()
	ex.subscribers = append(ex.subscribers, ch)

	// Send initial snapshot immediately
	statuses := make([]JobStatus, 0, len(ex.status))
	for _, s := range ex.status {
		statuses = append(statuses, *s)
	}
	ex.mu.Unlock()

	ch <- statuses
	return ch
}
