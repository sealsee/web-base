package setting

type UploadFile struct {
	Type       string `mapstructure:"type"`
	DomainName string `mapstructure:"domain_name"`
	*S3        `mapstructure:"s3"`
	*OSS       `mapstructure:"oss"`
	*Localhost `mapstructure:"localhost"`
}

type S3 struct {
	AccessKeyId     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"Secret_access_key"`
	Region          string `mapstructure:"region"`
	BucketName      string `mapstructure:"bucket_name"`
}

type OSS struct {
	AccessKeyId     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"Secret_access_key"`
	Endpoint        string `mapstructure:"end_point"` //VPC下用内部地址
	PublicBucket    string `mapstructure:"public_bucket"`
	PublicDomain    string `mapstructure:"public_domain"`
	PrivateBucket   string `mapstructure:"private_bucket"`
	PrivateDomain   string `mapstructure:"private_domain"`
}

type Localhost struct {
	PublicResourcePrefix  string `mapstructure:"public_resource_prefix"`
	PrivateResourcePrefix string `mapstructure:"private_resource_prefix"`
}
