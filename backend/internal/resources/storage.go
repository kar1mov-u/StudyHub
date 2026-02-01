package resources

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/joho/godotenv"
)

type FileStorage interface {
	UploadObject(ctx context.Context, filename string, size int64, body io.Reader) (string, error)
	DeleteObject(ctx context.Context, filename string) error
	CreatePresingedURL(ctx context.Context, key string) (string, error)
}

type S3Storage struct {
	s3Client   *s3.Client
	presigner  *s3.PresignClient
	bucketName string
	URL        string
}

func NewS3Storage(bucketname, url string) *S3Storage {
	_ = godotenv.Load()
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(client)
	return &S3Storage{s3Client: client, bucketName: bucketname, URL: url, presigner: presigner}
}

func (s *S3Storage) UploadObject(ctx context.Context, filename string, size int64, body io.Reader) (string, error) {
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(filename),
		Body:          body,
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", s.URL, filename), nil
}

func (s *S3Storage) DeleteObject(ctx context.Context, key string) error {
	// deleted := false
	input := s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}
	_, err := s.s3Client.DeleteObject(ctx, &input)
	if err != nil {
		var noKey *types.NoSuchKey
		var apiErr *smithy.GenericAPIError
		if errors.As(err, &noKey) {
			log.Printf("Object %s does not exist in %s.\n", key, s.bucketName)
			err = noKey
		} else if errors.As(err, &apiErr) {
			switch apiErr.ErrorCode() {
			case "AccessDenied":
				log.Printf("Access denied: cannot delete object %s from %s.\n", key, s.bucketName)
				err = nil
			}
		}
	}

	return err
}

func (s *S3Storage) CreatePresingedURL(ctx context.Context, key string) (string, error) {
	request, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(60 * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			s.bucketName, key, err)
	}
	return request.URL, err
}
