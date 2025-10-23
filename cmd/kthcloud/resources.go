package main

import "fmt"

const (
	kthcloud_cli = `   __    __    __         __                __             __   _ 
  / /__ / /_  / /  ____  / / ___  __ __ ___/ / ____ ____  / /  (_)
 /  '_// __/ / _ \/ __/ / / / _ \/ // // _  / /___// __/ / /  / / 
/_/\_\ \__/ /_//_/\__/ /_/  \___/\_,_/ \_,_/       \__/ /_/  /_/  
                                                                  `
	kthcloud_cli_blue = "\033[34m" + kthcloud_cli + "\033[0m"
)

var (
	// Semver version of the app
	version = "v0.0.1-dev"
	// Latest commit when compiling
	commit = "not-provided"

	banner = fmt.Sprintf("%s\nVersion:%-10s\t\t\t\tCommit:%-10s ", kthcloud_cli, version, commit)
)
