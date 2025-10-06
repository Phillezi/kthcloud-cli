package scheduler

func (d *SchedulerImpl) Status() <-chan []JobStatus {
	if d.executor == nil {
		ch := make(chan []JobStatus)
		close(ch)
		return ch
	}

	return d.executor.subscribe()
}
