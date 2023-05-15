package common

// 删除标记
const (
	Normal  = 1 // 正常的
	Deleted = 2 // 已删除的
	Yes     = 1 // 是
	No      = 2 // 否
)

const (
	EXCEL int = 1
	DBF   int = 2
)

const (
	OK      = 1
	Disable = 2
)

var UserCodeMsgMap = map[int]string{
	OK:      "正常",
	Disable: "停用",
}
