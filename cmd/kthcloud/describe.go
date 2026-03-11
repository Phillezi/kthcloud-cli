package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var describeCmd = &cobra.Command{
	Use: "describe",
}

func init() {
	rootCmd.AddCommand(describeCmd)

	describeCmd.PersistentFlags().StringP("output", "o", "table", "Output format: table, json, yaml")
	describeCmd.PersistentFlags().Bool("stats", false, "Print request stats (RTT, count, status) to stderr")

	viper.BindPFlags(describeCmd.PersistentFlags())
}
