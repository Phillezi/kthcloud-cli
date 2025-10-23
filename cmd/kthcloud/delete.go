package main

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use: "delete",
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
