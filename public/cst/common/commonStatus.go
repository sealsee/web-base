package common

const (
	OK      string = "0"
	Disable string = "1"
	Deleted string = "2"
)

const (
	EXCEL int = 1
	DBF   int = 2
)

var UserCodeMsgMap = map[string]string{
	OK:      "正常",
	Disable: "停用",
	Deleted: "删除",
}
