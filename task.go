package sysbench

import (
	"database/sql"
)

type TestTask interface {
	PrepareTask
	RunTask
}

type RunTask interface {
	LargeTxner
	SmallTxner
}

type PrepareTask interface {
	CreateTable(db *sql.DB) error
	InsertData(worker *Worker, db *sql.DB) error
}

type LargeTxner interface {
	LargeTxn(db *sql.DB) error
}

type SmallTxner interface {
	SmallTxn(worker *Worker, db *sql.DB) error
}
