package service

import (
	"context"
	"log"
	"os"

	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Service handles S3 file uploads.
type S3Service struct {
	client *s3.Client
}

// NewS3Service creates a new S3Service.
func NewS3Service() *S3Service {
	cfg := aws.Config{
		Region: "ap-northeast-1",
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	}

	client := s3.NewFromConfig(cfg)
	return &S3Service{client: client}
}

// Upload uploads a file to S3.
func (s *S3Service) Upload(bucketName, objectKey string, body io.Reader) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   body,
	})
	if err != nil {
		log.Printf("Error uploading to S3: %v", err)
		return err
	}
	return nil
}
