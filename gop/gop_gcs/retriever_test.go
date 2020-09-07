package gop_gcs

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gen/pkg/servers/gop"
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
	fakeGCS.s.CreateObject(fakestorage.Object{
		BucketName: r.AppConfig.CacheLocation,
		Name:       "github.com/anz-bank/sysl/tests/bananatree.sysl.pb.json@e78f4afc524ad8d1a1a4740779731d706b7b079b",
		Content:    []byte(bananatreepbjson),
	})
	r.downloader = fakeGCS.downloadInMem
	req := app.NewObject("github.com/anz-bank/sysl/tests/bananatree.sysl", "e78f4afc524ad8d1a1a4740779731d706b7b079b")
	require.NoError(t, r.Retrieve(req))
	banana := &gop.Object{Repo: "github.com/anz-bank/sysl", Version: "e78f4afc524ad8d1a1a4740779731d706b7b079b", Resource: "tests/bananatree.sysl", Content: "Bananatree [package=\"bananatree\"]:\n  !type Banana:\n    id \u003c: int\n    title \u003c: string\n\n  /banana:\n    /{id\u003c:int}:\n      GET:\n        return Banana\n\n  /morebanana:\n    /{id\u003c:int}:\n      GET:\n        return Banana\n"}
	req.Processed = nil
	require.Equal(t, banana, req)
}

func (s fakeServer) downloadInMem(bucket, object string) (io.Reader, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()
	return s.s.Client().Bucket(bucket).Object(object).NewReader(ctx)
}
