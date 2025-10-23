package aws

import (
	appConfig "github.com/raymondsugiarto/coffee-api/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config interface {
	LoadDefaultConfig() *s3.Client
}

type s3Config struct {
}

func NewS3Config() S3Config {
	return &s3Config{}
}

func (s *s3Config) LoadDefaultConfig() *s3.Client {
	awsConfig := appConfig.GetConfig().Aws
	accessKeyId := awsConfig.AccessKeyId
	secretKeyAccess := awsConfig.SecretKeyAccess

	options := s3.Options{
		Region:      "ap-southeast-1",
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyId, secretKeyAccess, "")),
	}

	return s3.New(options)

	// // Load the Shared AWS Configuration (~/.aws/config)
	// cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Create an Amazon S3 service client
	// return s3.NewFromConfig(cfg)

}
