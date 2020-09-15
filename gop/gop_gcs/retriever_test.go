package gop_gcs

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/joshcarp/gop/app"
	"github.com/stretchr/testify/require"
)

func TestGCSRetrieve(t *testing.T) {
	r := New(app.AppConfig{CacheLocation: "bucket"})
	fakeGCS := fakeServer{
		s: fakestorage.NewServer([]fakestorage.Object{}),
	}
	fakeGCS.s.CreateObject(fakestorage.Object{
		BucketName: r.AppConfig.CacheLocation,
		Name:       "github.com/anz-bank/sysl/tests/bananatree.sysl@e78f4afc524ad8d1a1a4740779731d706b7b079b",
		Content:    []byte(bananatree),
	})
	r.downloader = fakeGCS.downloadInMem
	obj, cached, err := r.Retrieve("github.com/anz-bank/sysl/tests/bananatree.sysl@e78f4afc524ad8d1a1a4740779731d706b7b079b")
	require.NoError(t, err)
	require.True(t, cached)
	require.Equal(t, []byte(bananatree), obj)
}

func (s fakeServer) downloadInMem(bucket, object string) (io.Reader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()
	return s.s.Client().Bucket(bucket).Object(object).NewReader(ctx)
}
