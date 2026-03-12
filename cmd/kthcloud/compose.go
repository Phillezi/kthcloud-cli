package main

import (
	"github.com/spf13/cobra"
)

var composeCmd = &cobra.Command{
	Use: "compose",
}

func init() {
	rootCmd.AddCommand(composeCmd)
}
