package export

import (
	"strconv"

	"github.com/sealsee/web-base/public/utils/export/internal"
	"github.com/xuri/excelize/v2"
)

func Dbf() internal.ImpExp {
	return nil
}

func Excel() internal.ImpExp {
	return &internal.Excel{}
}

// Deprecated
func ExportExcel(dataList [][]interface{}) (data []byte) {
	f := excelize.NewFile()
	defer f.Close()
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
