package IOFile

import (
	"bytes"
	"io"
	"mime/multipart"
	"path/filepath"
)

type fileParams struct {
	keyName     string
	contentType string
	data        io.Reader
	buf         *bytes.Buffer
}

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

func NewFileParamsRandomName(keyName string, file multipart.File) *fileParams {
	f := new(fileParams)
	f.keyName = keyName
	f.data = file
	f.contentType = contentType[filepath.Ext(keyName)]
	return f
}

func NewFileParamsNameBuffer(keyName string, buf *bytes.Buffer) *fileParams {
	f := new(fileParams)
	f.keyName = keyName
	f.buf = buf
	f.contentType = contentType[filepath.Ext(keyName)]
	return f
}
