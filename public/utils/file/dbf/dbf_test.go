package dbf

import (
	"fmt"
	"testing"

	"github.com/sealsee/web-base/public/utils/file/dbf/godbf"
)

func TestRowData(t *testing.T) {

	rows := make([][]string, 0, 3)
	row1 := []string{"code", "name", "state"}
	row2 := []string{"101", "技术", "正常"}
	row3 := []string{"102", "运营", "正常"}

	rows = append(rows, row1, row2, row3)

	// 写入dbf文件
	dbfTable := godbf.New("GBK")

	// 生成表头fields
	for _, field := range rows[0] {
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

	godbf.SaveToFile(dbfTable, "./test.dbf")

	readDbfTable, err := godbf.NewFromFile("./test.dbf", "GBK")
	if err != nil {
		fmt.Println(err)
	}

	// 读取dbf文件，取值验证
	dataList := make([][]string, readDbfTable.NumberOfRecords())

	for i := 0; i < readDbfTable.NumberOfRecords(); i++ {
		dataList[i] = make([]string, 3)
		dataList[i][0] = readDbfTable.FieldValue(i, 0)
		dataList[i][1] = readDbfTable.FieldValue(i, 1)
		dataList[i][2] = readDbfTable.FieldValue(i, 2)
	}

	fmt.Println(dataList)
}
