package tx

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init(d *sqlx.DB) {
	db = d
}

// Deprecated
type Tx interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...any) (sql.Result, error)
}

// Deprecated
func ExecTx(exec func(Tx)) bool {
	tx, err := db.Beginx()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			tx.Rollback()
			// panic(err)
		}
	}()
	exec(tx)
	tx.Commit()
	return true
}
