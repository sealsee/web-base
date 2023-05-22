package fileUtils

import (
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sealsee/web-base/public/utils/stringUtils"
)

var fileExtension = map[string]string{
	"image/png":  "png",
	"image/jpg":  "jpg",
	"image/jpeg": "jpeg",
	"image/bmp":  "bmp",
	"image/gif":  "gif",
}

var defaultAllowedExtension = []string{
	// 图片
	"bmp", "gif", "jpg", "jpeg", "png",
	// word excel powerpoint
	"doc", "docx", "xls", "xlsx", "ppt", "pptx", "html", "htm", "txt",
	// 压缩文件
	"rar", "zip", "gz", "bz2",
	// 视频格式
	"mp4", "avi", "rmvb",
	// pdf
	"pdf",
	// dbf
	"dbf",
}

func GetFileNameRandom(userId int64, extensionName string) string {
	t := time.Now()
	nameKey := t.Format("06/01/02") + "/" + stringUtils.GetUUID() + "_" + strconv.FormatInt(userId, 10) + "." + extensionName
	return nameKey
}

func GetExtension(file *multipart.FileHeader) string {
	ext := strings.Replace(filepath.Ext(file.Filename), ".", "", -1)
	if ext == "" {
		return fileExtension[file.Header["Content-Type"][0]]
	}
	return ext
}
