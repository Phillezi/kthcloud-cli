package main

import (
	"github.com/Phillezi/kthcloud-cli/cmd"
	"github.com/spf13/viper"
)

var buildTimestamp = "19700101000000"

func main() {
	viper.Set("release", "release-"+buildTimestamp)
	cmd.Execute()
}
