package cst

const (
	LoginUserKey             = "loginUser"
	LogKey                   = "log"
	ResourcePrefix           = "/profile" // 资源映射路径 前缀
	DefaultPublicPath        = "./file/public/"
	DefaultPrivatePath       = "./file/private/"
	UploadMaxNum             = 50               // 多文件上传最多支持50个文件
	UploadMaxSize      int64 = 10 * 1024 * 1024 // 单文件最大支持10MB
)

const (
	USER_TYPE_SYS = "00" // 系统分配用户（B端）
	USER_TYPE_REG = "01" // 注册用户（C端）
)
