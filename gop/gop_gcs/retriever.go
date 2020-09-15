package gop_gcs

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/joshcarp/gop/gop"

	"cloud.google.com/go/storage"
)

type downloader func(bucket, object string) (io.Reader, error)

func (a GOP) Retrieve(resource string) ([]byte, bool, error) {
	r, err := a.downloader(a.bucket, resource)
	if err != nil {
		return nil, false, gop.CreateError(gop.FileNotFoundError, "Error finding resource in cache", err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, false, gop.CreateError(gop.FileNotFoundError, "Error finding resource in cache", err)
	}
	return b, true, nil
}

func download(bucket, object string) (io.Reader, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()
	return client.Bucket(bucket).Object(object).NewReader(ctx)
}
