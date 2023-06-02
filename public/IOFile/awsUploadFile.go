package IOFile

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3IOFile struct {
	s3Config   *s3.Client
	bucket     string
	domainName string
}

func (l *s3IOFile) Upload(data io.Reader, suffixName, fileExt string, isPrivate bool) (string, error) {
	return "", nil
}

func (l *s3IOFile) Download(url string) ([]byte, error) {
	return nil, nil
}
