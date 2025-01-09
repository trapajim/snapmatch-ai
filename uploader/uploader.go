package uploader

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"google.golang.org/api/googleapi"
	"io"
)

type Uploader struct {
	client *storage.Client
	bucket string
}

func NewUploader(client *storage.Client, bucket string) *Uploader {
	return &Uploader{
		client: client,
		bucket: bucket,
	}
}

func (u *Uploader) Upload(ctx context.Context, file io.Reader, object string) error {
	wc := u.client.Bucket(u.bucket).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return handleApiError(err)
	}
	err := wc.Close()
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
	return snapmatchai.NewError(err, "error occurred, during file upload", 500)
}