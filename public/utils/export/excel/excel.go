package excel

import (
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func ExportExcel(rows [][]string) (data []byte) {
	//https://xuri.me/excelize/zh-hans/sheet.html#DuplicateRow
	// 创建一个工作表
	//index := f.NewSheet("Sheet1")
	//var a = make([]string, 0, 1)
	// 设置单元格的值
	//f.SetSheetRow("Sheet1", "1", &a)
	//
	//f.SetCellValue("Sheet1", "A2", a)
	// 设置工作簿的默认工作表
	//f.SetActiveSheet(index)

	f := excelize.NewFile()
	for i, row := range rows {
		f.SetSheetRow("Sheet1", "A"+strconv.Itoa(i+1), &row)
	}
	buffer, _ := f.WriteToBuffer()
	return buffer.Bytes()
}

type ExcelImpCall interface {
	Headers() []string
	Data() []any
}

type ExcelExpCall interface {
	Before()
	Extract() []any
	After()
}

func ImportExcel(bytes []byte, impCall ExcelImpCall) error {
	return nil
}

func ImportExcelWithUrl(ossUrl string, impCall ExcelImpCall) error {
	return nil
}

func ExportExcelSync(headers map[string]string, expCall ExcelExpCall) ([]byte, error) {
	return nil, nil
}

func ExportExcelAsync(headers map[string]string, expCall ExcelExpCall) (string, error) {
	return "", nil
}
