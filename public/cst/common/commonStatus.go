package common

// 删除标记
const (
	Normal  = 1 // 正常的
	Deleted = 2 // 已删除的
)

const (
	OPER_TYPE_IMPORT_DATA   int = 1 // 导入数据
	OPER_TYPE_IMPORT_UPDATE int = 2 // 导入更新数据
	OPER_TYPE_EXPORT        int = 3 // 导出
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
