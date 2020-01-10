package template

import (
	"testing"

	"github.com/pingcap/tidb-test/large_txn"
)

type MyPrepare struct {
	large.DefaultPrepareTask
}

// You can override the methods of the default one as you wish.
func (_ MyPrepare) CreateTable(db *sql.DB) error {
	sql1 := `create table if not exists sbtest1 (
id int(11) not null primary key,
k int(11) not null,
c char(120) not null default '',
pad char(255) not null default '')`
	_, err := db.Exec(sql1)
	if err != nil {
		return errors.WithStack(err)
	}
}

// Write your own test config, let it implement the TestTask interface.
type myTestTask struct {
	MyPrepare
	large.LargeDML
	large.SelectRandomRanges
}

// And run the test!
func TestT(t *testing.T) {
	// You can modify the test config.
	UserFlag = "root"
	HostFlag = "127.0.0.1"
	PortFlag = 6868
	DBFlag = "test"
	WorkerCount = 8

	task := myTestTask{}
	// Some of the builtin task also provide custom parameters that you can change.
	task.DefaultPrepareTask.InsertCount = 5000
	task.DefaultPrepareTask.RowsEachInsert = 50

	RunTest(task)
}
