package excel

import (
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func ExportExcel(dataList [][]interface{}) (data []byte) {
	f := excelize.NewFile()
	for i, row := range dataList {
		if i == 0 {
			f.SetSheetRow("Sheet1", "A1", &row)
		} else {
			f.SetSheetRow("Sheet1", "A"+strconv.Itoa(i+1), &row)
		}
	}
	buffer, _ := f.WriteToBuffer()
	return buffer.Bytes()
}

type ExcelImpCall interface {
	Headers() []string
	Row() map[string]interface{}
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
