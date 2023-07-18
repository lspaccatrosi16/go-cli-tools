package aws

import (
	"bytes"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	S3Client *s3.Client
	Bucket   string
}

const partMibs int64 = 10

func (b Bucket) UploadFile(key string, file []byte) {
	buf := bytes.NewReader(file)

	uploader := manager.NewUploader(b.S3Client, func(u *manager.Uploader) {
		u.PartSize = partMibs * 1024 * 1024
	})

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(key),
		Body:   buf,
	})

	if err != nil {
		log.Fatalln(err)
	}
}

func (b Bucket) GetFile(key string) []byte {

	downloader := manager.NewDownloader(b.S3Client, func(d *manager.Downloader) {
		d.PartSize = partMibs * 1024 * 1024
	})

	buffer := manager.NewWriteAtBuffer([]byte{})

	_, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		log.Fatalln(err)
	}

	return buffer.Bytes()
}

func NewBucket(sdkConfig aws.Config, bucketName string) Bucket {
	s3Client := s3.NewFromConfig(sdkConfig)
	bucket := Bucket{S3Client: s3Client, Bucket: bucketName}
	return bucket
}
