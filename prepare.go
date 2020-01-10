package sysbench

import (
	"bytes"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pingcap/errors"
	"github.com/spf13/cobra"
)

func getDSN() string {
	return fmt.Sprintf("%s@tcp(%s:%d)/%s", UserFlag, HostFlag, PortFlag, DBFlag)
}

type Worker struct {
	ID    int
	Count int
	succ  int
	fail  int
	dur   []float64
}

func CmdPrepare(cmd *cobra.Command, args []string) {
	db, err := sql.Open("mysql", getDSN())
	handleErr(err)
	db.SetMaxOpenConns(512)

	task := DefaultPrepareTask{}
	prepareTask(task, db)
}

func prepareTask(task PrepareTask, db *sql.DB) error {
	err := task.CreateTable(db)
	handleErr(err)

	var wg sync.WaitGroup
	wg.Add(WorkerCount)
	errs := make([]error, WorkerCount)
	for id := 0; id < WorkerCount; id++ {
		worker := Worker{ID: id, Count: WorkerCount, dur: make([]float64, 0, 100)}
		go worker.prepare(task, db, id, &wg, &errs[id])
	}
	wg.Wait()

	var lastErr error
	for i, err := range errs {
		if err != nil {
			fmt.Printf("worker %d error = %s\n", i, err.Error())
			lastErr = err
		}
	}
	fmt.Println("cmd prepare finish", errs)
	return lastErr
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func nextPrimaryID(workerCount int, current int) int {
	return current + workerCount
}

func (worker *Worker) prepare(task PrepareTask, db *sql.DB, workerID int, wg *sync.WaitGroup, res *error) {
	defer wg.Done()
	fmt.Printf("start prepare worker %d\n", workerID)
	err := task.InsertData(worker, db)
	if err != nil {
		*res = errors.WithStack(err)
		return
	}
	fmt.Printf("worker %d finish\n", workerID)
}

const ascii = "abcdefghijklmnopqrstuvwxyz1234567890"

func randString(n int) string {
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		pos := rand.Intn(len(ascii))
		buf.WriteByte(ascii[pos])
	}
	return buf.String()
}

type DefaultPrepareTask struct {
	InsertCount    int
	RowsEachInsert int
}

func (t DefaultPrepareTask) insertCount() int {
	if t.InsertCount > 0 {
		return t.InsertCount
	}
	return 1000
}

func (t DefaultPrepareTask) rowsEachInsert() int {
	if t.RowsEachInsert > 0 {
		return t.RowsEachInsert
	}
	return 30
}

func (_ DefaultPrepareTask) CreateTable(db *sql.DB) error {
	sql1 := `create table if not exists sbtest1 (
id int(11) not null primary key,
k int(11) not null,
c char(120) not null default '',
pad char(255) not null default '')`
	_, err := db.Exec(sql1)
	if err != nil {
		return errors.WithStack(err)
	}

	sql2 := `create table if not exists sbtest2 (
id int(11) not null,
k int(11) not null,
c char(120) not null default '',
pad char(255) not null default '')`
	_, err = db.Exec(sql2)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (task DefaultPrepareTask) InsertData(worker *Worker, db *sql.DB) error {
	var buf bytes.Buffer
	pkID := worker.ID
	for i := 0; i < task.insertCount(); i++ {
		buf.Reset()
		buf.WriteString("insert into sbtest1 (id, k, c, pad) values ")
		for i := 0; i < task.rowsEachInsert(); i++ {
			pkID = nextPrimaryID(worker.Count, pkID)
			dot := ""
			if i > 0 {
				dot = ", "
			}
			fmt.Fprintf(&buf, "%s(%d, %d, '%s', '%s')", dot, pkID, rand.Intn(1<<11), randString(32), randString(32))
		}

		_, err := db.Exec(buf.String())
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
