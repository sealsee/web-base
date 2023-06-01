package dbf

import "github.com/sealsee/web-base/public/utils/file/dbf/godbf"

func SetRows(rows [][]string) (data []byte) {

	dbfTable := godbf.New("GBK")

	fields := rows[0]

	// 生成表头fields
	for _, field := range fields {
		dbfTable.AddTextField(field, 32)
	}

	for i := 1; i < len(rows); i++ {
		// 先创建一行记录
		rowId, _ := dbfTable.AddNewRecord()
		// 遍历一行的每个值，放到field里
		for f, val := range rows[i] {
			dbfTable.SetFieldValue(rowId, f, val)
		}
	}

	return godbf.GetDbfFileData(dbfTable)
}

func ExportDBF(rows [][]string) (data []byte) {
	dbfTable := godbf.New("GBK")

	fields := rows[0]

	// 生成表头fields
	for _, field := range fields {
		dbfTable.AddTextField(field, 32)
	}

	for i := 1; i < len(rows); i++ {
		// 先创建一行记录
		rowId, _ := dbfTable.AddNewRecord()
		// 遍历一行的每个值，放到field里
		for f, val := range rows[i] {
			dbfTable.SetFieldValue(rowId, f, val)
		}
	}

	return godbf.GetDbfFileData(dbfTable)
}
