package uploader

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/trapajim/snapmatch-ai/server/middleware"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"google.golang.org/api/googleapi"
	"io"
	"log"
	"strings"
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
	wc := u.client.Bucket(u.defaultBucket).Object(objectName(ctx, object)).NewWriter(ctx)
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
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(expiry),
	}
	signedURL, err := u.client.Bucket(u.defaultBucket).SignedURL(objectName(ctx, object), opts)
	if err != nil {
		return "", handleApiError(err)
	}
	return signedURL, nil
}

func (u *Uploader) GetFile(ctx context.Context, object string) (io.ReadCloser, error) {
	rc, err := u.client.Bucket(u.defaultBucket).Object(objectName(ctx, object)).NewReader(ctx)
	if err != nil {
		return nil, handleApiError(err)
	}
	return rc, nil
}

func (u *Uploader) UpdateMetadata(ctx context.Context, object string, metadata map[string]string) error {
	_, err := u.client.Bucket(u.defaultBucket).Object(objectName(ctx, object)).Update(ctx, storage.ObjectAttrsToUpdate{
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

func objectName(ctx context.Context, object string) string {
	sess := middleware.GetSession(ctx)
	if sess == nil {
		log.Println("session is nil")
		return object
	}
	if strings.Contains(object, sess.SessionID()) {
		log.Println("object already contains session id")
		return object
	}
	log.Println("object does not contain session id")
	log.Println(fmt.Sprintf("%s/%s", sess.SessionID(), object))
	return fmt.Sprintf("%s/%s", sess.SessionID(), object)
}
