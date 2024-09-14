package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sealsee/web-base/public/basemodel"
	"github.com/sealsee/web-base/public/ds/page"
	"github.com/sealsee/web-base/public/utils/jsonUtils"
	"gorm.io/gorm"
)

var gormdb *gorm.DB

func InitGTx(gdb *gorm.DB) {
	gormdb = gdb
}

// 处理where条件, 转换合并成map条件+自定义条件
func convertWhereQuery(where basemodel.IQuery) (map[string]interface{}, string, []interface{}) {
	// 获取设置的表别名
	alias := where.GetAlias()
	if alias != "" {
		alias += "."
	}
	columns, conditions, args := where.GetConditions()
	whereMap, _ := jsonUtils.StructToDbMap(where)
	for k, v := range whereMap {
		// 删除旧key
		delete(whereMap, k)
		var hasCol bool
		// 判断condition column是否在where条件里，如果包含则map里去除
		for _, col := range columns {
			if col == k {
				hasCol = true
			}
		}
		if !hasCol && k != "curPage" && k != "pageSize" {
			// 添加新key
			whereMap[alias+k] = v
		}
	}
	condStr := ""
	for i := 1; i <= len(conditions); i++ {
		suffix := "AND"
		if i == len(conditions) {
			suffix = ""
		}
		condStr += fmt.Sprintf("%v%v %v ", alias, conditions[i-1], suffix)
	}
	return whereMap, condStr, args
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
	if where == nil {
		return 0
	}
	whereMap, conditions, condArgs := convertWhereQuery(where)
	t := new(T)
	var count int64
	gdb := gormdb.Model(t).Where(whereMap)
	if conditions != "" {
		gdb = gdb.Where(conditions, condArgs...)
	}
	if query != nil {
		gdb = gdb.Where(query, args...)
	}
	rlt := gdb.Count(&count)
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
	if query != nil {
		gdb = gdb.Where(query, args...)
	}
	rlt := gdb.Order(where.GetOrders()).Find(&ts)

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
	gdb := gormdb.Model(t).Offset(page.GetOffset()).Limit(page.GetLimit()).Where(whereMap)
	if conditions != "" {
		gdb = gdb.Where(conditions, condArgs...)
	}
	if query != nil {
		gdb = gdb.Where(query, args...)
	}
	rlt := gdb.Order(where.GetOrders()).Find(&ts)

	if rlt.RowsAffected <= 0 {
		return nil
	}
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return ts
}

// 查询指定字段, 返回字段的map列表
func ExecQueryListWithColumns[T any](columns []string, where basemodel.IQuery, page *page.Page, query interface{}, args ...interface{}) []map[string]any {
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
	gdb := gormdb.Model(t).Select(columns).Offset(page.GetOffset()).Limit(page.GetLimit()).Where(whereMap)
	if conditions != "" {
		gdb = gdb.Where(conditions, condArgs...)
	}
	if query != nil {
		gdb = gdb.Where(query, args...)
	}
	rlt := gdb.Order(where.GetOrders()).Find(&ts)

	if rlt.RowsAffected <= 0 {
		return nil
	}
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return ts
}

// 原生sql查询
func RawSqlQuery[T any](sql string, args ...interface{}) T {
	var ts T
	if sql == "" {
		return ts
	}
	rlt := gormdb.Raw(sql, args...).Scan(&ts)
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return ts
}

// 原生sql查询
func RawSqlQueryList[T any](sql string, args ...interface{}) (res []*T) {
	if sql == "" {
		return nil
	}
	ts := []*T{}
	rlt := gormdb.Raw(sql, args...).Scan(&ts)
	if rlt.RowsAffected <= 0 {
		return nil
	}
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	return ts
}

// 原生sql查询，支持列表分页
func RawSqlQueryListWithPage[T any](page *page.Page, sql string, args ...interface{}) (res []*T) {
	if page == nil || sql == "" {
		return nil
	}
	ts := []*T{}
	var total int64
	rlt := gormdb.Table("("+sql+") AS CT", args...).Count(&total)
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	sql += fmt.Sprintf(" LIMIT %v OFFSET %v", page.GetLimit(), page.GetOffset())
	rlt = gormdb.Raw(sql, args...).Scan(&ts)
	if rlt.RowsAffected <= 0 {
		return nil
	}
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	page.SetTotalSize(int(total))
	return ts
}

// 原生sql查询，支持列表分页，支持自定义条件where condition
func RawSqlQueryListWithPageWhere[T any](where basemodel.IQuery, page *page.Page, sql string, args ...interface{}) (res []*T) {
	if page == nil {
		return nil
	}
	ts := []*T{}
	var total int64

	whereMap, conditions, condArgs := convertWhereQuery(where)
	for k, v := range whereMap {
		sql += fmt.Sprintf(" AND %s = '%v'", k, v)
	}
	if conditions != "" {
		sql += fmt.Sprintf(" AND %v", conditions)
		args = append(args, condArgs...)
	}

	rlt := gormdb.Table("("+sql+") AS CT", args...).Count(&total)
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	sql += fmt.Sprintf(" LIMIT %v OFFSET %v", page.GetLimit(), page.GetOffset())
	rlt = gormdb.Raw(sql, args...).Scan(&ts)
	if rlt.RowsAffected <= 0 {
		return nil
	}
	if rlt.Error != nil {
		panic(rlt.Error)
	}
	page.SetTotalSize(int(total))
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
