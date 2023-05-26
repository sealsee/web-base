package IOFile

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type s3IOFile struct {
	s3Config   *s3.Client
	bucket     string
	domainName string
}

func (s *s3IOFile) PublicUploadFile(file *fileParams) (string, error) {
	obj := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(file.keyName),
		Body:        file.data,
		ContentType: aws.String(file.contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	}
	_, err := s.s3Config.PutObject(context.TODO(), obj)
	if err != nil {
		return "", err
	}
	return s.domainName + "/" + file.keyName, nil
}

func (s *s3IOFile) PrivateUploadFile(file *fileParams) (string, error) {
	obj := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String("private/" + file.keyName),
		Body:        file.data,
		ContentType: aws.String(file.contentType),
		ACL:         types.ObjectCannedACLPrivate,
	}
	_, err := s.s3Config.PutObject(context.TODO(), obj)
	if err != nil {
		return "", err
	}
	return file.keyName, nil
}

func (l *s3IOFile) GetFileFullName(filename string) (string, error) {
	// TODO
	return "", nil
}

func (l *s3IOFile) Upload(data io.Reader, suffixName, fileExt string, isPrivate bool) (string, error) {
	return "", nil
}

func (l *s3IOFile) Download(url string) ([]byte, error) {
	return nil, nil
}
