package scheduler

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Sched struct {
	jobs map[string]*Job
	// maps a job to the jobs that depend on it
	depMap     map[*Job][]*Job
	depMu      sync.Mutex
	ctx        context.Context
	jobctx     context.Context
	cancel     context.CancelFunc
	mu         sync.Mutex
	runnableQ  chan *Job
	resultChan chan *Job
	revertOnce sync.Once
	disabled   bool
}

func NewSched(ctx context.Context) *Sched {
	jobctx, cancel := context.WithCancel(ctx)
	return &Sched{
		jobs:       make(map[string]*Job),
		depMap:     make(map[*Job][]*Job),
		ctx:        ctx,
		jobctx:     jobctx,
		cancel:     cancel,
		runnableQ:  make(chan *Job),
		resultChan: make(chan *Job),
	}
}

// blocking function
func (s *Sched) Start() {
	for {
		select {
		case <-s.ctx.Done():
			// stop everything and wait
			logrus.Debugln("scheduler cancelled")
			return
		case runnable := <-s.runnableQ:
			logrus.Debugln("new runnable job")
			s.startJob(runnable, s.resultChan)
		case job := <-s.resultChan:
			if job != nil {
				go s.handleJobResult(job)
			}
		}
	}
}

func (s *Sched) revertJobs() {
	s.revertOnce.Do(func() {
		s.disabled = true
		for _, j := range s.jobs {
			switch j.State {
			case Started:
				j.mu.Lock()
				j.State = Cancelling
				j.mu.Unlock()
				j.cancel()
				j.CancelCallback()
				j.mu.Lock()
				j.State = Cancelled
				j.mu.Unlock()
			case Created:
				j.mu.Lock()
				j.State = Cancelled
				j.mu.Unlock()
			case Done:
				j.mu.Lock()
				j.State = Cancelling
				j.mu.Unlock()
				j.CancelCallback()
				j.mu.Lock()
				j.State = Cancelled
				j.mu.Unlock()
			}
		}
	})
}

func (s *Sched) getRunnable() *Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range s.jobs {
		if v != nil && v.State == Created {
			hasUnfinishedDeps := false
			for _, dep := range v.Dependencies {
				if dep.State != Done {
					// job has a dep that isnt done
					hasUnfinishedDeps = true
					break
				}
			}
			if !hasUnfinishedDeps {
				return v
			}
		}
	}
	logrus.Debugln("no jobs available to be scheduled")
	return nil
}

func (s *Sched) handleJobResult(job *Job) {
	switch job.State {
	case Done:
		logrus.Debugln("job with id: " + job.ID + " is done")
		for _, v := range s.depMap[job] {
			if v != nil && v.State == Created {
				hasUnfinishedDeps := false
				if v.Dependencies != nil {
					for _, dep := range v.Dependencies {
						if dep != nil && dep.State != Done {
							// job has a dep that isnt done
							hasUnfinishedDeps = true
							break
						}
					}
				}
				if !hasUnfinishedDeps {
					logrus.Debugln("found ready job")
					s.runnableQ <- v
					logrus.Debugln("sent")
				} else {
					logrus.Debugln("unable to find ready job")
				}
			}
		}
	case Errored:

	case Cancelled:
		logrus.Debugln("job with id: " + job.ID + " was cancelled")
		for _, v := range s.depMap[job] {
			if v != nil && v.cancel != nil {
				v.mu.Lock()
				v.State = Cancelling
				v.mu.Unlock()
				v.CancelCallback()
				v.mu.Lock()
				v.State = Cancelled
				v.mu.Unlock()
			}
		}
	default:
		logrus.Warnln("Unexpected state to handle, ", job.State)
	}
}

func (s *Sched) updateRunnable() {
	job := s.getRunnable()
	if job != nil {
		s.runnableQ <- job
	}
}

func (s *Sched) queueIfRunnable(job *Job) {
	if job.Dependencies != nil {
		for _, dep := range job.Dependencies {
			if dep != nil && dep.State != Done {
				return
			}
		}
	}
	logrus.Debugln("queueing job ", job.ID)
	s.runnableQ <- job
}

func (s *Sched) startJob(runnable *Job, onDone chan *Job) {
	if runnable.State != Created {
		logrus.Debugln("tried to start job that cant be started")
		return
	}
	if s.disabled {
		return
	}
	runnable.mu.Lock()
	runnable.State = Started
	runnable.mu.Unlock()

	logrus.Debugln("Starting job with id: " + runnable.ID)

	go func(job *Job, onDone chan *Job) {
		defer func() {
			job.mu.Lock()
			if job.State == Started {
				job.State = Done
			} else if job.State == Cancelling {
				job.State = Cancelled
			}
			job.mu.Unlock()
			onDone <- job
		}()

		if err := job.Action(job.ctx, func() {
			s.revertJobs()
			logrus.Debugln("Cancelling job with id: " + runnable.ID)
			job.mu.Lock()
			job.State = Cancelling
			job.mu.Unlock()
			job.CancelCallback()
		}); err != nil {
			logrus.Debugln("Error occurred on job with id: " + runnable.ID)
			job.mu.Lock()
			job.State = Errored
			job.mu.Unlock()
			s.revertJobs()
		}
	}(runnable, onDone)
}

func (s *Sched) AddJob(job *Job) string {
	id := uuid.New().String()

	for _, ok := s.jobs[id]; ok; {
		id = uuid.New().String()
	}

	job.mu.Lock()
	job.ID = id
	job.ctx, job.cancel = context.WithCancel(s.jobctx)
	job.mu.Unlock()

	s.mu.Lock()
	s.jobs[job.ID] = job
	s.mu.Unlock()

	s.depMu.Lock()
	for _, dep := range job.Dependencies {
		if _, ok := s.depMap[dep]; !ok {
			s.depMap[dep] = []*Job{job}
		} else {
			s.depMap[dep] = append(s.depMap[dep], job)
		}
	}
	s.depMu.Unlock()

	logrus.Debugln("Added job with id: " + job.ID)

	s.queueIfRunnable(job)

	return id
}

func (s *Sched) CancelJobsBlock() {
	s.cancel()
	<-s.jobctx.Done()
}

func (s *Sched) CancelJob(id string) error {
	s.mu.Lock()
	job, exists := s.jobs[id]
	s.mu.Unlock()
	if !exists {
		return fmt.Errorf("job with ID %s does not exist", id)
	}

	if job.State == Cancelled || job.State == Errored {
		logrus.Debugln("job already cancelled or errored")
		return nil
	}

	if job.State == Started {
		job.cancel()
	}

	job.mu.Lock()
	job.State = Cancelled
	job.mu.Unlock()

	for _, v := range s.depMap[job] {
		if v != nil && v.cancel != nil {
			v.cancel()
			v.mu.Lock()
			v.State = Cancelled
			v.mu.Unlock()
		}
	}

	logrus.Debugln("Cancelled job with id: " + job.ID)

	return nil
}

func (s *Sched) GetJobState(id string) (JobState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.jobs[id]
	if !exists {
		return 0, fmt.Errorf("job with ID %s does not exist", id)
	}

	return job.State, nil
}
