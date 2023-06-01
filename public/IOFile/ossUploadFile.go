package IOFile

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
)

type ossIOFile struct {
	client        *oss.Client
	publicBucket  string
	publicDomain  string
	privateBucket string
	privateDomain string
}

func (oss *ossIOFile) Upload(data io.Reader, suffixName, fileExt string, isPrivate bool) (string, error) {
	if data == nil {
		return "", nil
	}

	domain := oss.publicDomain
	bucketName := oss.publicBucket
	if isPrivate {
		bucketName = oss.privateBucket
		domain = oss.privateDomain
	}

	bucket, err := oss.client.Bucket(bucketName)
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}

	t := time.Now()
	var pathBuilder strings.Builder
	pathBuilder.WriteString(t.Format("2006-01-02"))
	pathBuilder.WriteString("/")

	fileName := generalFileName(suffixName, fileExt)
	pathBuilder.WriteString(fileName)

	err = bucket.PutObject(pathBuilder.String(), data)
	if err != nil {
		zap.L().Error(err.Error())
		return "", err
	}

	if domain != "" && domain[len(domain)-1] != '/' {
		domain += "/"
	}

	return domain + pathBuilder.String(), nil
}

func (oss *ossIOFile) Download(url string) ([]byte, error) {
	if url == "" && !strings.HasPrefix(url, "http") {
		return nil, nil
	}

	objectKey := ""
	bucketName := ""
	if strings.HasPrefix(url, oss.publicDomain) {
		bucketName = oss.publicBucket
		objectKey = url[len(oss.publicDomain):]
	} else if strings.HasPrefix(url, oss.privateDomain) {
		bucketName = oss.privateBucket
		objectKey = url[len(oss.publicDomain):]
	}

	if bucketName == "" || objectKey == "" {
		return nil, errors.New(url + " is error")
	}

	if objectKey[0] == '/' {
		objectKey = objectKey[1:]
	}

	bucket, err := oss.client.Bucket(bucketName)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	rs, err := bucket.GetObject(objectKey)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(rs)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return buf.Bytes(), nil
}

func (oss *ossIOFile) PublicUploadFile(file *FileParams) (string, error)  { return "", nil }
func (oss *ossIOFile) PrivateUploadFile(file *FileParams) (string, error) { return "", nil }
func (oss *ossIOFile) GetFileFullName(filename string) (string, error)    { return "", nil }
