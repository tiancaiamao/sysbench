package sysbench

import (
	"database/sql"
	"github.com/pingcap/errors"
	"math/rand"

	_ "github.com/go-sql-driver/mysql"
)

type SelectRandomPoints struct{}

func (_ SelectRandomPoints) Execute(worker *Worker, db *sql.DB) error {
	rows, err := db.Query("select id, k, c, pad from sbtest2 where k in (?, ?, ?)", rand.Intn(100000), rand.Intn(100000), rand.Intn(100000))
	defer rows.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type UpdateRandomPoints struct{}

func (_ UpdateRandomPoints) Execute(worker *Worker, db *sql.DB) error {
	rows, err := db.Query("update sbtest2 set k = k + 1 where id in (?, ?, ?)", rand.Intn(100000), rand.Intn(100000), rand.Intn(100000))
	defer rows.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type SelectRandomRanges struct{}

func (_ SelectRandomRanges) Execute(worker *Worker, db *sql.DB) error {
	// db.Query("SELECT count(k) FROM sbteste1 WHERE k BETWEEN ? AND ? OR k BETWEEN ? AND ?")
	return nil
}
