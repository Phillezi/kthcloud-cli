package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var composeCmd = &cobra.Command{
	Use: "compose",
}

func init() {
	rootCmd.AddCommand(composeCmd)

	composeCmd.PersistentFlags().StringArrayP("file", "f", nil, "Compose configuration files")
	viper.BindPFlag("compose.file", composeCmd.PersistentFlags().Lookup("file"))

	composeCmd.PersistentFlags().StringP("project-name", "p", "", "Project name")
	viper.BindPFlag("compose.projecname", composeCmd.PersistentFlags().Lookup("project-name"))
}
