package query

import (
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sealsee/web-base/public/ds/page"
	"gorm.io/gorm"
)

var gormdb *gorm.DB

func InitGTx(gdb *gorm.DB) {
	gormdb = gdb
}

func ExecGetQueryCount[QT any, T any](where *QT) int {
	ts := []*T{}
	var count int64
	gormdb.Model(ts).Where(where).Count(&count)
	return int(count)
}

func ExecQueryList[QT any, T any](where *QT, page *page.Page) []*T {
	ts := []*T{}
	var count int64
	gormdb.Model(ts).Where(where).Count(&count)
	if count < 1 {
		return nil
	}

	page.SetTotalSize(int(count))
	rlt := gormdb.Offset(page.GetOffset()).Limit(page.GetLimit()).Where(where).Find(&ts)
	if rlt.RowsAffected <= 0 {
		return nil
	}
	return ts
}

func HandListPageQuery[T any](db *sqlx.DB, query string, args interface{}, page *page.Page) (list []*T) {
	//sql := strings.ToUpper(query)
	countRow, err := db.NamedQuery("SELECT COUNT(*) "+query[strings.Index(query, " FROM "):], args)
	if err != nil {
		panic(err)
	}
	total := new(int)
	if countRow.Next() {
		countRow.Scan(total)
	}
	page.SetTotalSize(*total)
	list = make([]*T, 0, page.PageSize)

	if *total < 1 {
		return
	}

	query += " limit " + strconv.Itoa(page.GetOffset()) + "," + strconv.Itoa(page.GetLimit())
	listRows, err := db.NamedQuery(query, args)
	if err != nil {
		panic(err)
	}
	for listRows.Next() {
		data := new(T)
		err = listRows.StructScan(data)
		if err != nil {
			panic(err)
		}
		list = append(list, data)
	}
	defer listRows.Close()

	return list
}
