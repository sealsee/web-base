package cst

const (
	ResourcePrefix           = "resource" // 资源映射路径 前缀
	PublicTag                = "s1"
	PrivateTag               = "s2"
	DefaultPublicPath        = "./file/public/"
	DefaultPrivatePath       = "./file/private/"
	UploadMaxNum             = 50               // 多文件上传最多支持50个文件
	UploadMaxSize      int64 = 10 * 1024 * 1024 // 单文件最大支持10MB
)
