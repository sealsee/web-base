package IOFile

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/sealsee/web-base/public/IOFile/cst"
	"github.com/sealsee/web-base/public/utils/fileUtils"
	"go.uber.org/zap"
)

type localHostIOFile struct {
	publicPath  string
	privatePath string
	domainName  string
}

func (l *localHostIOFile) PublicUploadFile(file *FileParams) (string, error) {

	return "", nil
}

func (l *localHostIOFile) PrivateUploadFile(file *FileParams) (string, error) {
	return "", nil
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

	var pathAll strings.Builder
	var filePath strings.Builder
	filePath.WriteString(cst.ResourcePrefix)
	filePath.WriteString("/")

	if isPrivate {
		pathAll.WriteString(l.privatePath)
		if l.privatePath[len(l.privatePath)-1] != '/' {
			pathAll.WriteString("/")
		}

		filePath.WriteString(cst.PrivateTag)
		filePath.WriteString("/")
	} else {
		pathAll.WriteString(l.publicPath)
		if l.publicPath[len(l.publicPath)-1] != '/' {
			pathAll.WriteString("/")
		}

		filePath.WriteString(cst.PublicTag)
		filePath.WriteString("/")
	}

	t := time.Now()
	filePath.WriteString(t.Format("2006-01-02"))
	filePath.WriteString("/")
	fileName := GeneralFileName(suffixName, fileExt)

	pathAll.WriteString(filePath.String())
	err = fileUtils.CreateMutiDir(pathAll.String())
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}

	pathAll.WriteString(fileName)
	file, err := os.OpenFile(pathAll.String(), os.O_CREATE|os.O_WRONLY, os.ModePerm)
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

	writer.Flush()
	filePath.WriteString(fileName)

	domain := l.domainName
	if domain != "" && domain[len(domain)-1] != '/' {
		l.domainName += "/"
	}

	return domain + filePath.String(), nil
}

func (l *localHostIOFile) Download(url string) ([]byte, error) {
	if url == "" && !strings.HasPrefix(url, "http") {
		return nil, nil
	}

	filePath := ""
	if strings.HasPrefix(url, l.domainName) {
		filePath = url[len(l.domainName):]
	}

	if filePath == "" {
		return nil, errors.New(url + " is error")
	}

	var pathBulider strings.Builder
	if idx := strings.Index(filePath, cst.ResourcePrefix+"/"+cst.PublicTag); idx != -1 {
		pathBulider.WriteString(l.publicPath)
	} else if idx := strings.Index(filePath, cst.ResourcePrefix+"/"+cst.PrivateTag); idx != -1 {
		pathBulider.WriteString(l.privatePath)
	} else {
		return nil, errors.New(url + " is error")
	}

	pathBulider.WriteString(filePath)
	if !fileUtils.IsExist(pathBulider.String()) {
		return nil, errors.New(pathBulider.String() + " is not exist")
	}
	bytes, err := ioutil.ReadFile(pathBulider.String())
	if err != nil {
		zap.L().Error(err.Error())
	}
	return bytes, err
}
