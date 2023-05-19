package query

import (
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sealsee/web-base/public/basemodel"
	"github.com/sealsee/web-base/public/ds/page"
	"gorm.io/gorm"
)

var gormdb *gorm.DB

func InitGTx(gdb *gorm.DB) {
	gormdb = gdb
}

func ExecGetQueryCount[QT any, T any](where *QT) int {
	if where == nil {
		return 0
	}
	t := new(T)
	var count int64
	rlt := gormdb.Model(t).Where(where).Count(&count)
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(count)
}

func ExecQueryList[QT any, T any](where basemodel.IQuery, page *page.Page) []*T {
	if where == nil || page == nil {
		return nil
	}

	ts := []*T{}
	var count int64
	rlt := gormdb.Model(ts).Where(where).Count(&count)
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	if count < 1 {
		return nil
	}

	page.SetTotalSize(int(count))
	rlt = gormdb.Offset(page.GetOffset()).Limit(page.GetLimit()).Where(where).Order(where.GetOrders()).Find(&ts)
	if rlt.RowsAffected <= 0 {
		return nil
	}
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return ts
}

func ExecGetQueryCountWithCondition[T any](where basemodel.IQuery, query interface{}, args ...interface{}) int {
	if (where == nil && query == nil) || args == nil {
		return 0
	}

	t := new(T)
	var count int64
	rlt := gormdb.Model(t).Where(where).Where(query, args...).Count(&count)
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(count)
}

func ExecQueryListWithCondition[T any](where basemodel.IQuery, page *page.Page, query interface{}, args ...interface{}) []*T {
	if page == nil {
		return nil
	}

	t := new(T)
	var count int64
	rlt := gormdb.Model(t).Where(where).Where(query, args...).Count(&count)
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	if count < 1 {
		return nil
	}

	ts := []*T{}
	page.SetTotalSize(int(count))
	rlt = gormdb.Offset(page.GetOffset()).Limit(page.GetLimit()).Where(where).Where(query, args...).Order(where.GetOrders()).Find(&ts)
	if rlt.RowsAffected <= 0 {
		return nil
	}
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return ts
}

func ExecQueryListMapWithCondition[T any](where basemodel.IQuery, page *page.Page, query interface{}, args ...interface{}) []map[string]any {
	if page == nil {
		return nil
	}

	t := new(T)
	var count int64
	rlt := gormdb.Model(t).Where(where).Where(query, args...).Count(&count)
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	if count < 1 {
		return nil
	}

	var ts []map[string]any
	page.SetTotalSize(int(count))
	rlt = gormdb.Model(t).Offset(page.GetOffset()).Limit(page.GetLimit()).Where(where).Where(query, args...).Order(where.GetOrders()).Find(&ts)
	if rlt.RowsAffected <= 0 {
		return nil
	}
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return ts
}

func ExecQueryListWithColumns[T any](columns []string, where basemodel.IQuery, query interface{}, args ...interface{}) []*T {
	if columns == nil || where == nil {
		return nil
	}

	ts := []*T{}

	rlt := gormdb.Select(columns).Where(where).Where(query, args...).Limit(100).Find(&ts)
	if rlt.RowsAffected <= 0 {
		return nil
	}
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return ts
}

// 废弃
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
