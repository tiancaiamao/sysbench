package sysbench

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/pingcap/errors"
)

type InsertIntoSelect struct{}

func (_ InsertIntoSelect) LargeTxn(db *sql.DB) error {
	fmt.Println("run large txn InsertIntoSelect")
	_, err := db.Exec("insert into sbtest2 select * from sbtest1")
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println("run large txn InsertIntoSelect finish")
	return nil
}

type LargeDML struct{}

func (conf LargeDML) LargeTxn(db *sql.DB) error {
	fmt.Println("run large txn DML")
	tx, err := db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	for i := 0; i < 50000; i++ {
		_, err := tx.Exec("update sbtest1 set k = k + 1 where id in (?, ?, ?)", rand.Intn(100000), rand.Intn(100000), rand.Intn(100000))
		if err != nil {
			return errors.WithStack(err)
		}
		if i%1000 == 0 {
			fmt.Println("current executed = ", i)
		}
	}
	err = tx.Commit()
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println("run large txn DML finish")
	return nil
}

type LargeUpdate struct{}

func (_ LargeUpdate) LargeTxn(db *sql.DB) error {
	fmt.Println("run large txn UPDATE")
	_, err := db.Exec("update sbtest1 set k = k + 1")
	fmt.Println("run large txn UPDATE finish!")
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type LargeDelete struct{}

func (_ LargeDelete) LargeTxn(db *sql.DB) error {
	fmt.Println("run large txn DELETE")
	_, err := db.Exec("delete from sbtest1")
	fmt.Println("run large txn DELETE finish!")
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

type LoadData struct{}

func (_ LoadData) LargeTxn(db *sql.DB) error {
	// db.Exec("load data local into sbtest2")
	fmt.Println("run large txn LoadData")
	time.Sleep(20 * time.Second)
	return nil
}
