package gop_gcs

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/joshcarp/gop/gop"

	"cloud.google.com/go/storage"
)

type uploader func(bucket string, object string, r io.Reader) error

func (a GOP) Cache(resource string, content []byte) (err error) {
	if err := a.upload(a.bucket, resource, bytes.NewReader(content)); err != nil {
		return fmt.Errorf("%s: %w", gop.CacheWriteError, err)
	}
	return nil
}

func UploadFile(bucket string, object string, r io.Reader) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	// Open local file.
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}
