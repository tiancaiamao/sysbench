package sysbench

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/VividCortex/gohistogram"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

func CmdRun(cmd *cobra.Command, args []string) {
	db, err := sql.Open("mysql", getDSN())
	handleErr(err)
	db.SetMaxOpenConns(512)

	task := DefaultRunTask{}
	runTask(task, db)
}

func runTask(task RunTask, db *sql.DB) error {
	notify := make(chan struct{})

	report := &Report{
		Hist: gohistogram.NewHistogram(160),
	}
	exit := make(chan struct{})
	sampleCh := make(chan []float64, 10)
	go func(hist *gohistogram.NumericHistogram, input chan []float64, total int, exit chan struct{}) {
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
	}(report.Hist, sampleCh, WorkerCount, exit)

	workers := make([]Worker, WorkerCount)
	for i := 0; i < WorkerCount; i++ {
		workers[i] = Worker{ID: i, Count: WorkerCount, dur: make([]float64, 0, 100)}
		go workers[i].run(task, db, sampleCh, notify)
	}

	start := time.Now()
	err := largeTxn(task, db, notify)
	report.Duration = time.Since(start)
	if err != nil {
		fmt.Println("run large txn fail:", err)
	}
	<-exit

	for i := 0; i < WorkerCount; i++ {
		report.Succ += workers[i].succ
		report.Fail += workers[i].fail
	}

	report.Report()

	return err
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
		err := task.SmallTxn(worker, db)
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

type DefaultRunTask struct {
	LargeUpdate
	SelectRandomPoints
	UpdateRandomPoints
}

func largeTxn(task LargeTxner, db *sql.DB, notify chan struct{}) error {
	defer close(notify)
	err := task.LargeTxn(db)
	if err != nil {
		return err
	}
	return nil
}

func (t DefaultRunTask) SmallTxn(worker *Worker, db *sql.DB) error {
	t.SelectRandomPoints.SmallTxn(worker, db)
	t.UpdateRandomPoints.SmallTxn(worker, db)
	return nil
}

func RunTest(task TestTask) {
	db, err := sql.Open("mysql", getDSN())
	handleErr(err)

	err = prepareTask(task, db)
	handleErr(err)

	err = runTask(task, db)
	handleErr(err)
}
