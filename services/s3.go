package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"os"
)

type S3Serivce struct {
	s3Client *s3.Client
}

type S3Details struct {
	BucketName string
	Key        string
	FileName   string
}

func (s S3Serivce) UploadFile(bucketName, key, fileName string) error {
	log.Println("Upload DB")
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	return err
}

func (s S3Serivce) DownloadFile(bucketName, key, fileName string) error {
	log.Println("Download DB")
	s3Obj, err := s.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	defer s3Obj.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(s3Obj.Body)
	if err != nil {
		return err
	}
	_, err = file.Write(body)
	return err
}

func NewS3Service(config aws.Config) *S3Serivce {
	return &S3Serivce{
		s3.NewFromConfig(config),
	}
}
