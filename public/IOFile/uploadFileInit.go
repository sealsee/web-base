package IOFile

import (
	"fmt"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sealsee/web-base/public/IOFile/cst"
	"github.com/sealsee/web-base/public/setting"
	"github.com/sealsee/web-base/public/utils/set"
	"go.uber.org/zap"
)

const (
	awsS3     = "s3"
	osS       = "oss"
	localhost = "localhost"
)

var FileType = set.Set[string]{}

type IOFile interface {
	//Deprecated
	PublicUploadFile(file *FileParams) (string, error)
	//Deprecated
	PrivateUploadFile(file *FileParams) (string, error)
	//Deprecated
	GetFileFullName(filename string) (string, error)

	// fileExt, suffixName should be null
	Upload(data io.Reader, suffixName, fileExt string, isPrivate bool) (string, error)
	Download(url string) ([]byte, error)
}

var ioFile IOFile

func GetConfig() IOFile {
	return ioFile
}

func Init() {
	FileType.Add(awsS3)

	switch setting.Conf.UploadFile.Type {
	case awsS3:
		config := aws.Config{
			Credentials: credentials.NewStaticCredentialsProvider(setting.Conf.UploadFile.S3.AccessKeyId, setting.Conf.UploadFile.S3.SecretAccessKey, ""),
			Region:      setting.Conf.UploadFile.S3.Region,
		}

		s := new(s3IOFile)
		s.s3Config = s3.NewFromConfig(config)
		s.bucket = setting.Conf.UploadFile.S3.BucketName
		s.domainName = setting.Conf.UploadFile.DomainName
		ioFile = s
	case osS:
		client, err := oss.New(setting.Conf.UploadFile.OSS.Endpoint, setting.Conf.UploadFile.OSS.AccessKeyId, setting.Conf.UploadFile.OSS.SecretAccessKey)
		if err != nil {
			zap.L().Error(err.Error())
		}
		o := new(ossIOFile)
		o.client = client
		o.publicBucket = setting.Conf.UploadFile.OSS.PublicBucket
		o.publicDomain = setting.Conf.UploadFile.OSS.PublicDomain
		o.privateBucket = setting.Conf.UploadFile.OSS.PrivateBucket
		o.privateDomain = setting.Conf.UploadFile.OSS.PrivateDomain
		ioFile = o
	default:
		l := new(localHostIOFile)
		l.domainName = setting.Conf.UploadFile.DomainName
		pubPath := setting.Conf.UploadFile.Localhost.PublicResourcePrefix

		if setting.Conf.IsDocker || pubPath == "" {
			pubPath = cst.DefaultPublicPath
		}
		l.publicPath = pubPath

		priPath := setting.Conf.UploadFile.Localhost.PrivateResourcePrefix
		if setting.Conf.IsDocker || priPath == "" {
			priPath = cst.DefaultPrivatePath
		}
		l.privatePath = priPath
		ioFile = l
	}

	zap.L().Info(fmt.Sprintf("FileStore[%s] init success...\n", setting.Conf.UploadFile.Type))
}
