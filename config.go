package sysbench

import (
	"fmt"
	"time"
)

type Config struct {
	Conn    ConnConfig
	Prepare PrepareConfig
	Run     RunConfig
	Cleanup CleanupTask
}

func DefaultConfig() Config {
	return Config{
		Conn:    DefaultConnConfig(),
		Prepare: DefaultPrepareConfig(),
		Run:     DefaultRunConfig(),
		Cleanup: DefaultCleanupTask(),
	}
}

// Cmd is provided for the commandline flag.
var Cmd = DefaultConfig()

type ConnConfig struct {
	User string
	Host string
	Port int
	DB   string
}

func (conf *ConnConfig) getDSN() string {
	return fmt.Sprintf("%s@tcp(%s:%d)/%s", conf.User, conf.Host, conf.Port, conf.DB)
}

func DefaultConnConfig() ConnConfig {
	return ConnConfig{
		User: "root",
		Host: "127.0.0.1",
		Port: 4000,
		DB:   "test",
	}
}

type PrepareConfig struct {
	WorkerCount int
	Task        PrepareTask
}

func DefaultPrepareConfig() PrepareConfig {
	return PrepareConfig{
		WorkerCount: 4,
		Task:        DefaultPrepareTask(),
	}
}

type RunConfig struct {
	WorkerCount int
	Task        RunTask
	Duration    time.Duration
}

func DefaultRunConfig() RunConfig {
	return RunConfig{
		WorkerCount: 4,
		Task:        DefaultRunTask(),
		Duration:    time.Minute,
	}
}
