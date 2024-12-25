package spinner

import (
	"sync"

	"github.com/briandowns/spinner"
	"github.com/sirupsen/logrus"
)

type SpinnerHook struct {
	spinner *spinner.Spinner
	mu      sync.Mutex
}

func NewSpinnerHook(spinner *spinner.Spinner) *SpinnerHook {
	return &SpinnerHook{
		spinner: spinner,
	}
}

func (h *SpinnerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *SpinnerHook) Fire(entry *logrus.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.spinner != nil {
		h.spinner.Stop()
	}

	if h.spinner != nil && h.spinner.Enabled() {
		h.spinner.Start()
		/*go func() {
			time.Sleep(100 * time.Millisecond)
			h.spinner.Start()
		}()*/
	}

	return nil
}
