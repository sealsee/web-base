package excel

import (
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sealsee/web-base/public/utils/export/internal"
)

type Excel struct {
}

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

func (e *Excel) Import(bytes []byte, handler internal.ImpHandler) error {
	return nil
}

func (e *Excel) ImportWithUrl(url string, handler internal.ImpHandler) error {
	return nil
}

func (e *Excel) ExportSync(headers map[string]string, handler internal.ExpHandler) ([]byte, error) {
	return nil, nil
}

func (e *Excel) ExportAsync(headers map[string]string, handler internal.ExpHandler) (string, error) {
	return "", nil
}
