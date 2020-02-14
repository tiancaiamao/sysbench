package sysbench

import (
	"database/sql"
)

type dummpCleanup struct {
}

func (dummpCleanup) DropTable(db *sql.DB) error {
	return nil
}

func DefaultCleanupTask() CleanupTask {
	return dummpCleanup{}
}

func Cleanup(conf *Config) {
	if conf.Cleanup == nil {
		return
	}

	db, err := sql.Open("mysql", conf.Conn.getDSN())
	handleErr(err)
	defer db.Close()

	err = conf.Cleanup.DropTable(db)
	handleErr(err)
}
