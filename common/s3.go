package common

import (
	"distriai-index-solana/common/s3actions"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// initS3 initializes the S3 client and presigner with the necessary AWS configuration.
func initS3() {
	s3Client := s3.New(s3.Options{
		Region:      Conf.Aws.Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(Conf.Aws.AccessKeyId, Conf.Aws.SecretAccessKey, "")),
	})
	presignClient := s3.NewPresignClient(s3Client)
	S3Presigner = s3actions.Presigner{PresignClient: presignClient}
}
