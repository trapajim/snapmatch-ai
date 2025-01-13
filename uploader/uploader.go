package uploader

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"google.golang.org/api/googleapi"
	"io"
	"time"
)

type Uploader struct {
	client        *storage.Client
	defaultBucket string
}

func NewUploader(client *storage.Client, bucket string) *Uploader {
	return &Uploader{
		client:        client,
		defaultBucket: bucket,
	}
}
func (u *Uploader) WithBucket(bucket string) snapmatchai.Uploader {
	return &Uploader{
		client:        u.client,
		defaultBucket: bucket,
	}
}
func (u *Uploader) Upload(ctx context.Context, file io.Reader, object string) error {
	wc := u.client.Bucket(u.defaultBucket).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return handleApiError(err)
	}
	err := wc.Close()
	if err != nil {
		return handleApiError(err)
	}
	return nil
}

func (u *Uploader) SignUrl(ctx context.Context, object string, expiry time.Duration) (string, error) {
	u.client.Bucket(u.defaultBucket).Object(object)
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(expiry),
	}
	signedURL, err := u.client.Bucket(u.defaultBucket).SignedURL(object, opts)
	if err != nil {
		return "", handleApiError(err)
	}
	return signedURL, nil
}

func (u *Uploader) GetFile(ctx context.Context, object string) (io.ReadCloser, error) {
	rc, err := u.client.Bucket(u.defaultBucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, handleApiError(err)
	}
	return rc, nil
}

func (u *Uploader) UpdateMetadata(ctx context.Context, object string, metadata map[string]string) error {
	_, err := u.client.Bucket(u.defaultBucket).Object(object).Update(ctx, storage.ObjectAttrsToUpdate{
		Metadata: metadata,
	})
	if err != nil {
		return handleApiError(err)
	}
	return nil
}

func handleApiError(err error) error {
	var e *googleapi.Error
	if ok := errors.As(err, &e); ok {
		return snapmatchai.NewError(err, fmt.Sprintf("Google API error: %s, Code: %d", e.Message, e.Code), e.Code)
	}
	return snapmatchai.NewError(err, "error occurred, during file handling", 500)
}
