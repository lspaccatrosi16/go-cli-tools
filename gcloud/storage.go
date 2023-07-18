package gcloud

import (
	"bytes"
	"context"
	"io"
	"time"

	s2 "cloud.google.com/go/storage"
	s1 "firebase.google.com/go/storage"
)

type Bucket struct {
	Bucket *s2.BucketHandle
	ctx    context.Context
}

func (b *Bucket) GetFile(key string) []byte {
	buffer := bytes.NewBuffer([]byte{})

	rc, err := b.Bucket.Object(key).NewReader(b.ctx)

	if err != nil {
		panic(err)
	}

	io.Copy(buffer, rc)

	err = rc.Close()

	if err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (b *Bucket) UploadFile(key string, contents []byte) {
	buf := bytes.NewReader(contents)

	wc := b.Bucket.Object(key).NewWriter(b.ctx)

	io.Copy(wc, buf)

	err := wc.Close()

	if err != nil {
		panic(err)
	}
}

func (b *Bucket) GetTemporaryUrl(key string, expiry int) string {
	signOpts := s2.SignedURLOptions{
		Scheme:  s2.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(time.Duration(expiry) * time.Hour),
	}
	u, err := b.Bucket.SignedURL(key, &signOpts)

	if err != nil {
		panic(err)
	}

	return u
}

type GStorageClient struct {
	app    *FirebaseApp
	Client *s1.Client
}

func (s *GStorageClient) Close() {
	// theres no close function huh
}

func (s *GStorageClient) GetBucket(name string) *Bucket {
	bucket, err := s.Client.Bucket(name)

	if err != nil {
		panic(err)
	}

	return &Bucket{
		Bucket: bucket,
		ctx:    app.ctx,
	}
}

func NewGStorage() *GStorageClient {
	app := getFirebase()

	client, err := app.app.Storage(app.ctx)

	if err != nil {
		panic(err)
	}

	sClient := GStorageClient{
		app:    app,
		Client: client,
	}

	return &sClient
}
