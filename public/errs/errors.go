package errs

import "fmt"

type ERROR [2]string

// 是否有异常
func (e ERROR) Invalid() bool {
	return len(e) > 0
}

// 异常信息转为字符串，按msg(code)格式拼接
func (e ERROR) String() string {
	return fmt.Sprintf("%v(%v)", e[1], e[0])
}

// 异常内容变量赋值
func (e ERROR) Format(a ...any) ERROR {
	return ERROR{e[0], fmt.Sprintf(e[1], a...)}
}

var (
	UNAUTHORIZED = ERROR{"401", "认证失败"}

	REFUSE_VISIT_ERR     = ERROR{"000", "拒绝访问！！!"}
	USER_PASSWORD_ERR    = ERROR{"002", "用户不存在/密码错误"}
	USER_DELETED_ERR     = ERROR{"003", "对不起，您的账号：%s 已删除"}
	USER_DISABLED_ERR    = ERROR{"004", "对不起，您的账号：%s 已停用"}
	CAPTCHA_ERR          = ERROR{"005", "验证码错误"}
	PARAM_INVALID_ERR    = ERROR{"006", "参数错误"}
	HANDLE_ERR           = ERROR{"007", "操作失败"}
	FILE_UPLOAD_ERR      = ERROR{"008", "文件上传失败"}
	FILE_DOWNLOAD_ERR    = ERROR{"009", "文件下载失败"}
	FILE_UPLOAD_MAX_NUM  = ERROR{"010", "文件上传个数超出限制"}
	FILE_UPLOAD_MAX_SIZE = ERROR{"011", "文件上传大小超出限制"}
)
