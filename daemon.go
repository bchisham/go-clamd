package clamd

import (
	"os"
	"os/exec"
)

type daemonConfig struct {
	configFile *os.File
}

type DaemonOption func(*daemonConfig)

func WithConfigFile(configFile *os.File) DaemonOption {
	return func(c *daemonConfig) {
		c.configFile = configFile
	}
}

type Server struct {
	daemonConfig
}

func NewDaemon(opts ...DaemonOption) *Server {
	defaults := &daemonConfig{}
	for _, opt := range opts {
		opt(defaults)
	}
	return &Server{
		daemonConfig: *defaults,
	}
}

func (s *Server) Start() error {
	runner := exec.Command("clamd", "--config-file", s.configFile.Name())
	runner.Stdin = nil
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	return runner.Run()
}
