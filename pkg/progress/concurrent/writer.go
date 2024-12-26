package concurrent

import (
	"io"
	"sync"
)

type Writer struct {
	writer io.Writer   // Underlying writer to which logs are written.
	lock   *sync.Mutex // Mutex to synchronize access.
}

// Singleton instance of mutext.
var instance *sync.Mutex
var once sync.Once

// Write method for Writer, implements io.Writer interface.
// This ensures concurrent writes are safe.
func (cw *Writer) Write(p []byte) (n int, err error) {
	cw.lock.Lock()            // Lock to ensure thread safety.
	defer cw.lock.Unlock()    // Unlock after writing.
	return cw.writer.Write(p) // Write to the underlying writer.
}

func GetLock() *sync.Mutex {
	once.Do(func() {
		instance = &sync.Mutex{}
	})
	return instance
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
		lock:   GetLock(),
	}
}
