package errs

type ERROR [2]string

var REFUSE_VISIT_ERR = ERROR{"000", "拒绝访问！！!"}
var USER_UNAUTH_ERR = ERROR{"001", "认证失败！！!"}
var USER_PASSWORD_ERR = ERROR{"002", "用户不存在/密码错误"}
var USER_DELETED_ERR = ERROR{"003", "对不起，您的账号：%s 已删除"}
var USER_DISABLED_ERR = ERROR{"004", "对不起，您的账号：%s 已停用"}
var CAPTCHA_ERR = ERROR{"005", "验证码错误"}
var PARAM_INVALID_ERR = ERROR{"006", "参数错误"}
var HANDLE_ERR = ERROR{"007", "操作失败"}
var FILE_UPLOAD_ERR = ERROR{"008", "文件上传失败"}
var FILE_DOWNLOAD_ERR = ERROR{"009", "文件下载失败"}
