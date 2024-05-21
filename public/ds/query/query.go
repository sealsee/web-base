package query

import (
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sealsee/web-base/public/basemodel"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/utils/jsonUtils"
	"github.com/sealsee/web-base/public/utils/stringUtils"
	"gorm.io/gorm"
)

var gormdb *gorm.DB

func InitGTx(gdb *gorm.DB) {
	gormdb = gdb
}

// 处理where条件
func convertWhereQuery(where basemodel.IQuery) (map[string]interface{}, string, []interface{}) {
	columns, conditions, args := where.GetConditions()
	whereMap := jsonUtils.StructToMap(where)
	for k, v := range whereMap {
		dbCol := stringUtils.ToUnderScoreCase(k)
		// 删除旧key
		delete(whereMap, k)
		var hasCol bool
		// 判断condition column是否在where条件里，如果包含则map里去除
		for _, col := range columns {
			if col == dbCol {
				hasCol = true
			}
		}
		if !hasCol {
			// 添加新key
			whereMap[dbCol] = v
		}
	}
	return whereMap, conditions, args
}

func ExecGetQueryCount[QT any, T any](where basemodel.IQuery) int {
	if where == nil {
		return 0
	}
	whereMap, conditions, args := convertWhereQuery(where)
	t := new(T)
	var count int64

	gdb := gormdb.Model(t).Where(whereMap)
	if conditions != "" {
		gdb = gdb.Where(conditions, args...)
	}
	rlt := gdb.Count(&count)

	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(count)
}

func ExecQueryList[QT any, T any](where basemodel.IQuery, page *page.Page) []*T {
	if where == nil || page == nil {
		return nil
	}

	count := ExecGetQueryCount[QT, T](where)
	if count < 1 {
		return nil
	}

	whereMap, conditions, args := convertWhereQuery(where)

	page.SetTotalSize(count)
	gdb := gormdb.Offset(page.GetOffset()).Limit(page.GetLimit()).Where(whereMap)
	if conditions != "" {
		gdb = gdb.Where(conditions, args...)
	}
	ts := []*T{}
	rlt := gdb.Order(where.GetOrders()).Find(&ts)
	if rlt.RowsAffected <= 0 {
		return nil
	}
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return ts
}

// 特殊情况下使用，如需要拼接比较复杂的条件
func ExecGetQueryCountWithCondition[T any](where basemodel.IQuery, query interface{}, args ...interface{}) int {
	if (where == nil && query == nil) || args == nil {
		return 0
	}
	whereMap, conditions, condArgs := convertWhereQuery(where)
	t := new(T)
	var count int64
	gdb := gormdb.Model(t).Where(whereMap)
	if conditions != "" {
		gdb = gdb.Where(conditions, condArgs...)
	}
	rlt := gdb.Where(query, args...).Count(&count)
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return int(count)
}

func ExecQueryListWithCondition[T any](where basemodel.IQuery, page *page.Page, query interface{}, args ...interface{}) []*T {
	if page == nil {
		return nil
	}

	count := ExecGetQueryCountWithCondition[T](where, query, args...)
	if count < 1 {
		return nil
	}

	whereMap, conditions, condArgs := convertWhereQuery(where)

	ts := []*T{}
	page.SetTotalSize(count)
	gdb := gormdb.Offset(page.GetOffset()).Limit(page.GetLimit()).Where(whereMap)
	if conditions != "" {
		gdb = gdb.Where(conditions, condArgs...)
	}
	rlt := gdb.Where(query, args...).Order(where.GetOrders()).Find(&ts)

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

	count := ExecGetQueryCountWithCondition[T](where, query, args...)
	if count < 1 {
		return nil
	}

	whereMap, conditions, condArgs := convertWhereQuery(where)

	t := new(T)
	var ts []map[string]any
	page.SetTotalSize(count)
	gdb2 := gormdb.Model(t).Offset(page.GetOffset()).Limit(page.GetLimit()).Where(whereMap)
	if conditions != "" {
		gdb2 = gdb2.Where(conditions, condArgs...)
	}
	rlt := gdb2.Where(query, args...).Order(where.GetOrders()).Find(&ts)

	if rlt.RowsAffected <= 0 {
		return nil
	}
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return ts
}

// 查询指定字段，结果集暂时限制100条
func ExecQueryListWithColumns[T any](columns []string, where basemodel.IQuery, query interface{}, args ...interface{}) []*T {
	if columns == nil || where == nil {
		return nil
	}
	whereMap, conditions, condArgs := convertWhereQuery(where)
	ts := []*T{}

	// rlt := gormdb.Select(columns).Where(where).Where(query, args...).Limit(100).Find(&ts)
	gdb2 := gormdb.Select(columns).Where(whereMap)
	if conditions != "" {
		gdb2 = gdb2.Where(conditions, condArgs...)
	}
	rlt := gdb2.Where(query, args...).Limit(100).Find(&ts)

	if rlt.RowsAffected <= 0 {
		return nil
	}
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return ts
}

// Deprecated
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
