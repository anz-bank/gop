package cacher_gcs

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	retrieve "github.com/joshcarp/pb-mod/config"
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type Cacher struct {
	AppConfig retrieve.AppConfig
}

func (a Cacher) Cache(res *pbmod.Object) (err error) {
	filename := fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version)
	if err := UploadFile(a.AppConfig.CacheLocation, filename, strings.NewReader(res.Content)); err != nil {
		return err
	}
	filename = fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version)
	if err := UploadFile(a.AppConfig.CacheLocation, filename, strings.NewReader(*res.Processed)); err != nil {
		return err
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
