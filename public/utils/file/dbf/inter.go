package dbf

import "github.com/sealsee/web-base/public/ds/page"

type DBFFieldType int

const (
	S  DBFFieldType = iota + 1 //string
	N                          //number
	D                          //double
	DT                         //date
)

type DBFField struct {
	DBFname string //dbf 字段名
	Dname   string //数据字段
	Cname   string //中文名，可选
	Type    DBFFieldType
	Length  int
}

type ImpHandler interface {
	Headers([]string)
	Row(*map[string]string)
	After()
}

type ExpHandler interface {
	Title() string
	Fields() []DBFField
	Rows(*page.Page) []map[string]interface{}
	Finish(url string)
}

type ImpExp interface {
	GetHeaders(arg any) ([]string, error)
	Import(bytes []byte, handler ImpHandler) error
	ImportWithUrl(url string, handler ImpHandler) error
	Export(handler ExpHandler) ([]byte, error)
}
