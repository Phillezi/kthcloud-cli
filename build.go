//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

const (
	defaultTarget   = "all"
	defaultSelector = "./..."
	defaultOutdir   = "bin"
)

var (
	targets = map[string]string{
		defaultTarget: "Build everything",
	}
)

func main() {
	args := os.Args[1:]
	exitCode := 0

	if len(args) < 1 {
		args = append(args, defaultTarget)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	version := getGitVersion()

	for _, arg := range args {
		switch arg {
		case defaultTarget:
			if err := os.MkdirAll(defaultOutdir, os.ModePerm); err != nil {
				errf("Failed to create output directory %s: %v", defaultOutdir, err)
				os.Exit(1)
			}

			ldflags := fmt.Sprintf("-w -s -X main.version=%s", version)
			cmd := exec.CommandContext(ctx, "go",
				"build",
				"-ldflags", ldflags,
				"-o", defaultOutdir,
				defaultSelector,
			)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			infof("Building target '%s' with version %s", arg, version)
			if err := cmd.Run(); err != nil {
				errf("Error occurred when building: %v", err)
				exitCode = 1
			}
		default:
			errf("No target named: %s, available targets: %v", arg, keys(targets))
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

func getGitVersion() string {
	tag, _ := runGit("describe", "--tags", "--abbrev=0")
	commit, _ := runGit("rev-parse", "--short", "HEAD")

	// check if HEAD points to the tag
	tagCommit, _ := runGit("rev-list", "-n", "1", tag)
	if tagCommit == commit {
		return tag
	}
	return fmt.Sprintf("%s-dirty-%s", tag, commit)
}

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Stderr = nil
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
)

type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERROR
	DEBUG
)

func logf(level LogLevel, format string, args ...any) {
	var prefix, color string

	switch level {
	case INFO:
		prefix = "+"
		color = Green
	case WARN:
		prefix = "!"
		color = Yellow
	case ERROR:
		prefix = "*"
		color = Red
	case DEBUG:
		prefix = "#"
		color = Cyan
	default:
		prefix = "-"
		color = Reset
	}

	fmt.Fprintf(os.Stderr, "%s%s %s%s\n", color, prefix, fmt.Sprintf(format, args...), Reset)
}

func infof(format string, args ...any)  { logf(INFO, format, args...) }
func warnf(format string, args ...any)  { logf(WARN, format, args...) }
func errf(format string, args ...any)   { logf(ERROR, format, args...) }
func debugf(format string, args ...any) { logf(DEBUG, format, args...) }

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
