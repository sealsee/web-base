package IOFile

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/sealsee/web-base/public/cst"
	"github.com/sealsee/web-base/public/utils/fileUtils"
)

type localHostIOFile struct {
	publicPath  string
	privatePath string
	domainName  string
}

func (l *localHostIOFile) PublicUploadFile(file *fileParams) (string, error) {
	var b []byte
	if file.buf == nil {
		buf := &bytes.Buffer{}
		_, err := buf.ReadFrom(file.data)
		if err != nil {
			return "", err
		}
		b = buf.Bytes()
	} else {
		b = file.buf.Bytes()
	}
	pathAndName := l.publicPath + file.keyName
	err := fileUtils.CreateMutiDir(filepath.Dir(pathAndName))
	if err != nil {
		return "", err
	}
	err = os.WriteFile(pathAndName, b, 0664)
	if err != nil {
		return "", err
	}
	return cst.ResourcePrefix + "/" + file.keyName, nil
}

func (l *localHostIOFile) privateUploadFile(file *fileParams) (string, error) {
	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(file.data)
	if err != nil {
		return "", err
	}
	pathAndName := l.privatePath + file.keyName
	err = fileUtils.CreateMutiDir(filepath.Dir(pathAndName))
	if err != nil {
		return "", err
	}
	b := buf.Bytes()
	err = os.WriteFile(pathAndName, b, 0664)
	if err != nil {
		return "", err
	}
	return file.keyName, nil
}

func (l *localHostIOFile) GetFileFullName(filename string) (string, error) {
	if !strings.HasPrefix(filename, cst.ResourcePrefix+"/") {
		return "", errors.New("wrong path! should prefix with '/profile/'")
	}
	keyName := strings.Replace(filename, cst.ResourcePrefix+"/", "", 1)
	return l.publicPath + keyName, nil
}
