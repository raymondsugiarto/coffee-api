package aws

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type UploadFileInput struct {
	File       io.Reader
	BucketName string
	Key        string
}

type S3Service interface {
	UploadFile(input UploadFileInput) (*s3.PutObjectOutput, error)
	DeleteFile(input UploadFileInput) (*s3.DeleteObjectOutput, error)
}

type s3Service struct {
	s3Config S3Config
}

func NewS3Service(s3Config S3Config) S3Service {
	return &s3Service{
		s3Config: s3Config,
	}
}

func (s *s3Service) UploadFile(input UploadFileInput) (*s3.PutObjectOutput, error) {
	client := s.s3Config.LoadDefaultConfig()

	// Upload the file to S3
	output, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(input.BucketName),
		Key:    aws.String(input.Key),
		Body:   input.File,
		ACL:    types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		log.Fatalf("failed to upload file, %v", err)
		return nil, err
	}
	return output, err
}

func (s *s3Service) DeleteFile(input UploadFileInput) (*s3.DeleteObjectOutput, error) {
	client := s.s3Config.LoadDefaultConfig()

	// Delete file from S3
	output, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(input.BucketName),
		Key:    aws.String(input.Key),
	})

	if err != nil {
		log.Fatalf("failed to delete file, %v", err)
		return nil, err
	}
	return output, err
}
