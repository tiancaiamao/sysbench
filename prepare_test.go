package sysbench

import (
	"log"
	"testing"
)

func TestT(t *testing.T) {
	var conf Config
	conf.Conn = ConnConfig{
		User: "root",
		Host: "127.0.0.1",
		Port: 4000,
		DB:   "test",
	}
	conf.Prepare = PrepareConfig{
		WorkerCount: 4,
		Task:        DefaultPrepareTask(),
	}
	conf.Run = RunConfig{
		WorkerCount: 4,
		Task:        SelectRandomPoints{},
	}

	RunTest(&conf)
}
