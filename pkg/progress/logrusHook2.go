package progress

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type TrackerHook struct {
	tracker *Tracker
	mu      sync.Mutex
}

func NewTrackerHook(tracker *Tracker) *TrackerHook {
	return &TrackerHook{
		tracker: tracker,
	}
}

func (h *TrackerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *TrackerHook) Fire(entry *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	/*if h.tracker != nil {
		h.tracker.multi.Stop()
	}*/

	/*if h.tracker != nil && h.tracker.multi.IsActive {
		h.tracker.multi.Start()
	}*/

	return nil
}
