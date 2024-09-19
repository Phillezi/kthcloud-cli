package compose

import (
	"os"
	"os/signal"
	"syscall"
)

func setupSignalHandler(done chan bool, handler func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		<-sigs
		defer func() {
			done <- true
		}()

		handler()

	}()
}
