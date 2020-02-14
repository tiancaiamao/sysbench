package sysbench

import (
	"testing"
)

func TestT(t *testing.T) {
	conf := Config{
		Conn:    DefaultConnConfig(),
		Prepare: DefaultPrepareConfig(),
		Run: RunConfig{
			WorkerCount: 4,
			Task:        SelectRandomPoints{},
		},
	}

	RunTest(&conf)
}
