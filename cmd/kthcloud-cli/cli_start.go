package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
}
