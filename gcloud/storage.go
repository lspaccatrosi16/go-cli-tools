package gcloud

import (
	"bytes"
	"context"
	"io"
	"time"

	s2 "cloud.google.com/go/storage"
	s1 "firebase.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type Bucket struct {
	Bucket *s2.BucketHandle
	ctx    context.Context
}

func (b *Bucket) GetFile(key string) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	rc, err := b.Bucket.Object(key).NewReader(b.ctx)

	if err != nil {
		return []byte{}, wrapStorage(err)
	}

	io.Copy(buffer, rc)

	err = rc.Close()

	if err != nil {
		return []byte{}, wrapStorage(err)
	}

	return buffer.Bytes(), nil
}

func (b *Bucket) UploadFile(key string, contents []byte) error {
	buf := bytes.NewReader(contents)

	wc := b.Bucket.Object(key).NewWriter(b.ctx)

	io.Copy(wc, buf)

	err := wc.Close()

	if err != nil {
		return wrapStorage(err)
	}
	return nil
}

func (b *Bucket) DeleteFile(key string) error {
	obj := b.Bucket.Object(key)

	err := obj.Delete(b.ctx)

	if err != nil {
		return wrapStorage(err)
	}

	return nil
}

func (b *Bucket) ListKeys() ([]string, error) {
	keys := []string{}

	result := b.Bucket.Objects(b.ctx, nil)

	for {
		object, err := result.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return keys, wrapStorage(err)
		}

		keys = append(keys, object.Name)

	}

	return keys, nil
}

func (b *Bucket) GetTemporaryUrl(key string, expiry int) (string, error) {
	signOpts := s2.SignedURLOptions{
		Scheme:  s2.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(time.Duration(expiry) * time.Hour),
	}
	u, err := b.Bucket.SignedURL(key, &signOpts)

	if err != nil {
		return "", wrapStorage(err)
	}

	return u, nil
}

type GStorageClient struct {
	app    *FirebaseApp
	Client *s1.Client
}

func (s *GStorageClient) Close() {
	// theres no close function huh
}

func (s *GStorageClient) GetBucket(name string) (*Bucket, error) {
	bucket, err := s.Client.Bucket(name)

	if err != nil {
		return nil, wrapStorage(err)
	}

	return &Bucket{
		Bucket: bucket,
		ctx:    app.ctx,
	}, nil
}

func NewGStorage() (*GStorageClient, error) {
	app, err := getFirebase()

	if err != nil {
		return nil, wrapStorage(err)
	}

	client, err := app.app.Storage(app.ctx)

	if err != nil {
		return nil, wrapStorage(err)
	}

	sClient := GStorageClient{
		app:    app,
		Client: client,
	}

	return &sClient, nil
}
