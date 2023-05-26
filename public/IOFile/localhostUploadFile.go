package IOFile

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sealsee/web-base/public/cst"
	"github.com/sealsee/web-base/public/utils/fileUtils"
	"go.uber.org/zap"
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

func (l *localHostIOFile) PrivateUploadFile(file *fileParams) (string, error) {
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

func (l *localHostIOFile) Upload(data io.Reader, suffixName, fileExt string, isPrivate bool) (string, error) {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(data)
	if err != nil {
		return "", err
	}

	var pathBuilder strings.Builder
	if isPrivate {
		pathBuilder.WriteString(l.privatePath)
		if l.privatePath[len(l.privatePath)-1] != '/' {
			pathBuilder.WriteString("/")
		}
	} else {
		pathBuilder.WriteString(l.publicPath)
		if l.publicPath[len(l.publicPath)-1] != '/' {
			pathBuilder.WriteString("/")
		}
	}

	t := time.Now()
	pathBuilder.WriteString(t.Format("2006-01-02"))
	pathBuilder.WriteString("/")

	path := pathBuilder.String()
	fileName := GeneralFileName(suffixName, fileExt)
	pathBuilder.WriteString(fileName)
	filePath := pathBuilder.String()

	err = fileUtils.CreateMutiDir(path)
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(buf.Bytes())
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}

	return filePath, nil
}

func (l *localHostIOFile) Download(url string) ([]byte, error) {
	if url == "" {
		return nil, nil
	}

	relativePath := url
	if strings.HasPrefix(url, "http") {
		idx := strings.Index(url, l.domainName)
		if idx != -1 {
			relativePath = url[len(l.domainName):]
		}
	}

	if !fileUtils.IsExist(relativePath) {
		return nil, errors.New(url + " is not exist")
	}
	bytes, err := ioutil.ReadFile(relativePath)
	if err != nil {
		zap.L().Error(err.Error())
	}
	return bytes, err
}
