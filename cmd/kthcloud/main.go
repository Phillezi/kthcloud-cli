package main

import (
	"log"
	"os"

	"github.com/Phillezi/common/interrupt"
)

func main() {
	if err := rootCmd.ExecuteContext(interrupt.GetInstance().Context()); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
