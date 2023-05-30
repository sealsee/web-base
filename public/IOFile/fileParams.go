package IOFile

import (
	"bytes"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/sealsee/web-base/public/utils/stringUtils"
)

var contentType = map[string]string{
	".json": "application/json",
	".html": "text/html",
	".js":   "application/javascript",
	".css":  "text/css",
	".gif":  "image/gif",
	".png":  "image/png",
	".gz":   "application/x-gzip",
	".svg":  "image/svg+xml",
	".pdf":  "application/pdf",
	".jpeg": "image/jpeg",
}

// Deprecated
type FileParams struct {
	keyName     string
	contentType string
	data        io.Reader
	buf         *bytes.Buffer
}

// Deprecated
func NewFileParamsRandomName(keyName string, file multipart.File) *FileParams {
	f := new(FileParams)
	f.keyName = keyName
	f.data = file
	f.contentType = contentType[filepath.Ext(keyName)]
	return f
}

// Deprecated
func NewFileParamsNameBuffer(keyName string, buf *bytes.Buffer) *FileParams {
	f := new(FileParams)
	f.keyName = keyName
	f.buf = buf
	f.contentType = contentType[filepath.Ext(keyName)]
	return f
}

// Deprecated
func GeneralFileName(suffixName, fileExt string) string {
	var fileName strings.Builder
	fileName.WriteString(stringUtils.GetUUID())
	if suffixName != "" {
		fileName.WriteString("_")
		fileName.WriteString(suffixName)
	}
	if fileExt != "" {
		if fileExt[0] != '.' {
			fileName.WriteString(".")
		}
		fileName.WriteString(fileExt)
	}

	return fileName.String()
}
