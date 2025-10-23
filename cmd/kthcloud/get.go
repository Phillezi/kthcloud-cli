package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCmd = &cobra.Command{
	Use: "get",
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().BoolP("all", "a", false, "Get all")
	getCmd.PersistentFlags().String("by-user-id", "", "Get all by a userID")
	getCmd.PersistentFlags().StringP("output", "o", "table", "Output format: table, json, yaml")
	getCmd.PersistentFlags().Bool("stats", false, "Print request stats (RTT, count, status) to stderr")

	viper.BindPFlags(getCmd.PersistentFlags())
}
