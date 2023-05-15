package export

import (
	"github.com/sealsee/web-base/public/utils/export/excel"
	"github.com/sealsee/web-base/public/utils/export/internal"
)

func Dbf() internal.ImpExp {
	return nil
}

func Excel() internal.ImpExp {
	return &excel.Excel{}
}
