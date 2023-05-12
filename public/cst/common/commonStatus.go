package common

// 删除标记
const (
	Normal  = 1 // 正常的
	Deleted = 2 // 已删除的
)

const (
	OK      string = "0"
	Disable string = "1"
)

const (
	EXCEL int = 1
	DBF   int = 2
)

var UserCodeMsgMap = map[string]string{
	OK:      "正常",
	Disable: "停用",
}
