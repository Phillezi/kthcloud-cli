package main

import "github.com/spf13/cobra"

var createCmd = &cobra.Command{
	Use: "create",
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringP("output", "o", "table", "Output format: table, json, yaml")
}
