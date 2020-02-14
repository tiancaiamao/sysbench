package sysbench

import (
	"database/sql"
)

type RunTask interface {
	Execute(worker *Worker, db *sql.DB) error
}

type PrepareTask interface {
	CreateTable(db *sql.DB) error
	InsertData(worker *Worker, db *sql.DB) error
}

type CleanupTask interface {
	DropTable(db *sql.DB) error
}
