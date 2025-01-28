package clamd

import (
	"context"
	"os"
	"os/exec"
)

type freshclamConfig struct {
	nchecks int
	pidfile string
}

func WithNChecks(nchecks int) FreshclamOption {
	return func(c *freshclamConfig) {
		c.nchecks = nchecks
	}
}

func WithPidFile(pidfile string) FreshclamOption {
	return func(c *freshclamConfig) {
		c.pidfile = pidfile
	}
}

type FreshclamOption func(*freshclamConfig)

func Freshclam(ctx context.Context, opts ...FreshclamOption) error {
	defaults := &freshclamConfig{}
	for _, opt := range opts {
		opt(defaults)
	}

	c, cancel := context.WithCancel(ctx)
	defer cancel()
	var args []string

	if defaults.nchecks != 0 {
		args = append(args, "--checks="+string(defaults.nchecks))
	}
	if defaults.pidfile != "" {
		args = append(args, "--pidfile="+defaults.pidfile)
	}

	runner := exec.CommandContext(c, "freshclam", args...)
	runner.Stdin = nil
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	return runner.Run()
}
