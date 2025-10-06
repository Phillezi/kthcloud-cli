package scheduler

import "errors"

// Add a job and its dependencies
func (d *SchedulerImpl) Add(j Job, deps ...string) (id string, err error) {
	id, err = d.dag.AddVertex(j)
	if err != nil {
		return id, err
	}

	for _, dep := range deps {
		if e := d.dag.AddEdge(dep, id); e != nil {
			err = errors.Join(err, e)
		}
	}
	return
}
