package scheduler

import (
	"context"
	"log"
	"runtime"
	"sync"
)

// Run jobs concurrently when dependencies are satisfied
func (d *SchedulerImpl) Start() error {
	ctx, cancel := context.WithCancelCause(d.ctx)
	ex := &Executor{
		ctx:    ctx,
		cancel: cancel,
		done:   make(map[string]bool),
		status: make(map[string]*JobStatus),
	}

	// try to size jobCh to number of vertices; fall back to small buffer if unknown
	jobsCount := 0
	verts := d.dag.GetVertices()
	if verts != nil {
		jobsCount = len(verts)
	}
	if jobsCount <= 0 {
		jobsCount = 16
	}

	jobCh := make(chan string, jobsCount)
	var wg sync.WaitGroup

	workerCount := min(runtime.NumCPU(), jobsCount)
	if workerCount <= 0 {
		workerCount = 1
	}

	for id := range verts {
		ex.setState(id, JobPending, nil)
	}

	// start worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Go(func() {
			for id := range jobCh {
				// stop early if cancelled
				select {
				case <-ctx.Done():
					return
				default:
				}

				jobVal, err := d.dag.GetVertex(id)
				if err != nil {
					// unknown vertex (maybe removed) => skip
					continue
				}
				job, ok := jobVal.(Job)
				if !ok {
					// wrong type stored in DAG => skip
					continue
				}

				log.Println("Start:", id)
				ex.setState(id, JobRunning, nil)

				if err := job.Run(ctx); err != nil {
					log.Println("FAIL:", id, "Err:", err)

					ex.setState(id, JobFailed, err)
					ex.markFailed(id)

					// cancel the whole execution (stops other workers via context)
					cancel(err)
					return
				}
				log.Println("End:", id)
				ex.markDone(id)
				ex.setState(id, JobDone, nil)

				// Schedule dependents if they are ready
				children, err := d.dag.GetChildren(id)
				if err != nil {
					continue
				}
				for childID := range children {
					// only enqueue if all parents are done
					if ex.isReady(d, childID) {
						select {
						case jobCh <- childID:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		})
	}

	// enqueue initial root jobs (those without parents)
	roots := d.dag.GetRoots()
	for id := range roots {
		jobCh <- id
	}

	// close jobCh after all workers finished (wait in separate goroutine)
	go func() {
		wg.Wait()
		close(jobCh)
		for _, sub := range ex.subscribers {
			close(sub)
		}
	}()

	// wait until context is cancelled (by failure or external cancel)
	<-ctx.Done()

	// if a failure occurred, revert executed jobs in reverse order
	if ex.failed {
		log.Println("FAILURE detected => reverting executed jobs...")
		ex.revertAll(d)
	}

	return ctx.Err()
}
