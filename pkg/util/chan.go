package util

import "github.com/sirupsen/logrus"

func SafeSend[T any](ch chan T, data T) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Warn("Attempted to send to a closed channel")
		}
	}()
	ch <- data
}
