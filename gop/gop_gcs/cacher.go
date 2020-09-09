package gop_gcs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/joshcarp/gop/app"

	"cloud.google.com/go/storage"
	"github.com/joshcarp/gop/gop"
)

type uploader func(bucket string, object string, r io.Reader) error

func (a GOP) Cache(res gop.Object) (err error) {
	filename := fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version)
	if err := a.upload(a.AppConfig.CacheLocation, filename, bytes.NewReader(res.Content)); err != nil {
		return app.CreateError(app.CacheWriteError, "Error uploading file to cache", err)
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
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	wc := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	return nil
}
