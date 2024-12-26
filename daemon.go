package clamd

import (
	"context"
	"os"
	"os/exec"
)

type daemonConfig struct {
	configFile *os.File
	ctx        context.Context
}

type DaemonOption func(*daemonConfig)

func WithConfigFile(configFile *os.File) DaemonOption {
	return func(c *daemonConfig) {
		c.configFile = configFile
	}
}

func WithParentContext(ctx context.Context) DaemonOption {
	return func(c *daemonConfig) {
		c.ctx = ctx
	}
}

type Server struct {
	daemonConfig
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewDaemon(opts ...DaemonOption) *Server {
	defaults := &daemonConfig{}
	for _, opt := range opts {
		opt(defaults)
	}
	parentCtx := context.Background()
	if defaults.ctx != nil {
		parentCtx = defaults.ctx
	}
	c, cancel := context.WithCancel(parentCtx)
	return &Server{
		daemonConfig: *defaults,
		ctx:          c,
		cancelFunc:   cancel,
	}
}

func (s *Server) Start() error {
	runner := exec.CommandContext(s.ctx, "clamd", "--config-file", s.configFile.Name())
	runner.Stdin = nil
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	return runner.Run()
}

func (s *Server) Context() context.Context {
	return s.ctx
}

func (s *Server) Stop() {
	s.cancelFunc()
}
