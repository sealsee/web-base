package IOFile

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	"github.com/sealsee/web-base/public/cst"
	"github.com/sealsee/web-base/public/utils/fileUtils"
)

type localHostIOFile struct {
	publicPath  string
	privatePath string
	domainName  string
}

func (l *localHostIOFile) PublicUploadFile(file *fileParams) (string, error) {

	buf := &bytes.Buffer{}
	_, err := buf.ReadFrom(file.data)
	if err != nil {
		return "", err
	}
	b := buf.Bytes()
	pathAndName := l.publicPath + file.keyName
	err = fileUtils.CreateMutiDir(filepath.Dir(pathAndName))
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(pathAndName, b, 0664)
	if err != nil {
		return "", err
	}
	return l.domainName + cst.ResourcePrefix + "/" + file.keyName, nil
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
	err = ioutil.WriteFile(pathAndName, b, 0664)
	if err != nil {
		return "", err
	}
	return file.keyName, nil
}
