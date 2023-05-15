package common

// 删除标记
const (
	Normal  = 1 // 正常的
	Deleted = 2 // 已删除的
)

const (
	EXCEL int = 1
	DBF   int = 2
)

// 用户状态
const (
	OK      = "1"
	Disable = "0"
)

var UserCodeMsgMap = map[string]string{
	OK:      "正常",
	Disable: "停用",
}
