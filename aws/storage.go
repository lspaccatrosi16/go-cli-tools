package aws

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	S3Client *s3.Client
	Bucket   string
	Config   aws.Config
}

const partMibs int64 = 10

func (b *Bucket) UploadFile(key string, file []byte) error {
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
		return wrap(err)
	}

	return nil
}

func (b *Bucket) GetFile(key string) ([]byte, error) {

	downloader := manager.NewDownloader(b.S3Client, func(d *manager.Downloader) {
		d.PartSize = partMibs * 1024 * 1024
	})

	buffer := manager.NewWriteAtBuffer([]byte{})

	_, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return []byte{}, wrap(err)
	}

	return buffer.Bytes(), nil
}

func (b *Bucket) DeleteFile(key string) error {
	_, err := b.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return wrap(err)
	}

	return nil
}

func (b *Bucket) ListKeys() ([]string, error) {
	keys := []string{}

	result, err := b.S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(b.Bucket),
	})

	if err != nil {
		return keys, wrap(err)
	}

	for _, k := range result.Contents {
		keys = append(keys, *k.Key)
	}

	return keys, nil
}

func (b *Bucket) GetTemporaryUrl(key string, expiry int) (string, error) {
	presignClient := s3.NewPresignClient(b.S3Client)
	presignedUrl, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(b.Bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(time.Hour*time.Duration(expiry)))

	if err != nil {
		return "", wrap(err)
	}

	return presignedUrl.URL, nil
}

func (b *Bucket) GetObjectUrl(key string) string {
	region := b.Config.Region
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", b.Bucket, region, key)
}

func NewBucket(sdkConfig aws.Config, bucketName string) Bucket {
	s3Client := s3.NewFromConfig(sdkConfig)
	bucket := Bucket{S3Client: s3Client, Bucket: bucketName, Config: sdkConfig}
	return bucket
}
