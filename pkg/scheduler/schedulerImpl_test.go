package scheduler_test

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/kthcloud/cli/pkg/scheduler"
)

// --- Helper to get current goroutine ID ---
func getGID() uint64 {
	b := make([]byte, 64)
	n := runtime.Stack(b, false)
	b = b[:n]
	fields := strings.Fields(string(b))
	if len(fields) < 2 {
		return 0
	}
	id, _ := strconv.ParseUint(fields[1], 10, 64)
	return id
}

// --- Single log entry for thread-safe logging ---
type LogEntry struct {
	Msg string
	GID uint64
}

// --- Mock Job Implementation ---
type TestJob struct {
	Value string
	Fail  bool
	mu    *sync.Mutex
	logs  *[]LogEntry
}

func (j *TestJob) Run(ctx context.Context) error {
	gid := getGID()
	j.mu.Lock()
	*j.logs = append(*j.logs, LogEntry{Msg: fmt.Sprintf("start:%s", j.Value), GID: gid})
	j.mu.Unlock()

	time.Sleep(100 * time.Millisecond)

	if j.Fail {
		j.mu.Lock()
		*j.logs = append(*j.logs, LogEntry{Msg: fmt.Sprintf("fail:%s", j.Value), GID: gid})
		j.mu.Unlock()
		return fmt.Errorf("job %s failed", j.Value)
	}

	j.mu.Lock()
	*j.logs = append(*j.logs, LogEntry{Msg: fmt.Sprintf("done:%s", j.Value), GID: gid})
	j.mu.Unlock()
	return nil
}

func (j *TestJob) Revert(ctx context.Context) error {
	gid := getGID()
	j.mu.Lock()
	*j.logs = append(*j.logs, LogEntry{Msg: fmt.Sprintf("revert:%s", j.Value), GID: gid})
	j.mu.Unlock()
	return nil
}

// --- Test ---
func TestConcurrentDagExecution(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	dag := scheduler.New(ctx)

	mu := &sync.Mutex{}
	logs := &[]LogEntry{}

	// Graph: hello → world → foo → Bar(fails)
	helloID, _ := dag.Add(&TestJob{Value: "hello", mu: mu, logs: logs})
	worldID, _ := dag.Add(&TestJob{Value: "world", mu: mu, logs: logs}, helloID)
	fooID, _ := dag.Add(&TestJob{Value: "foo", mu: mu, logs: logs}, worldID)
	_, _ = dag.Add(&TestJob{Value: "Bar", Fail: true, mu: mu, logs: logs}, fooID)

	independantID, _ := dag.Add(&TestJob{Value: "independant", mu: mu, logs: logs})
	_, _ = dag.Add(&TestJob{Value: "dependant", mu: mu, logs: logs}, independantID)

	independant2ID, _ := dag.Add(&TestJob{Value: "independant2", mu: mu, logs: logs})
	_, _ = dag.Add(&TestJob{Value: "dependant2", mu: mu, logs: logs}, independant2ID)

	err := dag.Start()
	if err == nil {
		t.Errorf("expected DAG to fail but got no error")
	}

	// --- Assertions ---
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("Execution log:")
	for i, entry := range *logs {
		fmt.Printf("%d: %s (GID %d)\n", i, entry.Msg, entry.GID)
	}

	// Check concurrency
	concurrent := false
	gidSet := make(map[uint64]int)
	for _, entry := range *logs {
		gidSet[entry.GID]++
		if gidSet[entry.GID] > 1 {
			concurrent = true
			break
		}
	}

	if !concurrent {
		t.Errorf("expected concurrent execution, but jobs ran sequentially")
	}

	// Check failure and reverts
	foundFail := false
	foundRevert := false
	for _, entry := range *logs {
		if strings.HasPrefix(entry.Msg, "fail:") {
			foundFail = true
		}
		if strings.HasPrefix(entry.Msg, "revert:") {
			foundRevert = true
		}
	}

	if !foundFail {
		t.Errorf("expected a failing job but none failed")
	}
	if !foundRevert {
		t.Errorf("expected revert jobs but none were reverted")
	}
}
