package resources

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileStorage interface {
	UploadFile(ctx context.Context, filename string, body io.Reader)
}

type S3Storage struct {
	s3Client   *s3.Client
	bucketName string
}

func (s *S3Storage) UploadFile(ctx context.Context, filename string, body io.Reader) error {
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(filename),
		Body:   body,
	})
	return err
}
