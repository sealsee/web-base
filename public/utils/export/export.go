package export

import "github.com/sealsee/web-base/public/utils/export/excel"

func Dbf() excel.ImpExp {
	return nil
}

func Excel() excel.ImpExp {
	return &internal.Excel{}
}
