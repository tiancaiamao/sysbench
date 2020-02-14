package sysbench

import (
	"database/sql"
	// "fmt"
	// "log"
	// "net/http"
	_ "net/http/pprof"
	"time"

	"github.com/VividCortex/gohistogram"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/spf13/cobra"
)

func runTask(workerCount int, duration time.Duration, task RunTask, db *sql.DB) error {
	notify := make(chan struct{})

	report := &Report{
		Hist: gohistogram.NewHistogram(160),
	}
	exit := make(chan struct{})
	sampleCh := make(chan []float64, 10)
	go backgroundStatistics(report.Hist, sampleCh, workerCount, exit)

	workers := make([]Worker, workerCount)
	for i := 0; i < workerCount; i++ {
		workers[i] = Worker{ID: i, Count: workerCount, dur: make([]float64, 0, 100)}
		go workers[i].run(task, db, sampleCh, notify)
	}

	start := time.Now()
	time.Sleep(duration)
	close(notify)
	report.Duration = time.Since(start)
	<-exit

	for i := 0; i < workerCount; i++ {
		report.Succ += workers[i].succ
		report.Fail += workers[i].fail
	}

	report.Report()

	return nil
}

func (worker *Worker) run(task RunTask, db *sql.DB, sampleCh chan<- []float64, done <-chan struct{}) {
	for {
		select {
		case <-done:
			sampleCh <- worker.dur
			sampleCh <- nil
			return
		default:
		}

		start := time.Now()
		err := task.Execute(worker, db)
		elapse := time.Since(start)
		if err != nil {
			worker.fail++
		} else {
			worker.succ++
		}
		if len(worker.dur) == cap(worker.dur) {
			sz := cap(worker.dur)
			sampleCh <- worker.dur
			worker.dur = make([]float64, 0, sz)
		}
		worker.dur = append(worker.dur, elapse.Seconds()*1000)

	}
}

func backgroundStatistics(hist *gohistogram.NumericHistogram, input chan []float64, total int, exit chan struct{}) {
	finished := 0
	for data := range input {
		if data == nil {
			// use data == nil to mean finish
			finished++
			if finished == total {
				break
			}
		}
		for _, val := range data {
			hist.Add(val)
		}
	}
	close(exit)
}

func DefaultRunTask() RunTask {
	return baseRunTask{}
}

type baseRunTask struct{}

func (t baseRunTask) Execute(worker *Worker, db *sql.DB) error {
	return nil
}

func Run(conf *Config) {
	db, err := sql.Open("mysql", conf.Conn.getDSN())
	handleErr(err)
	defer db.Close()

	db.SetMaxOpenConns(512)

	err = runTask(conf.Run.WorkerCount, conf.Run.Duration, conf.Run.Task, db)
	handleErr(err)
}

func RunTest(conf *Config) {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	Prepare(conf)
	Run(conf)
	Cleanup(conf)
}
