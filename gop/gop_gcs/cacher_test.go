package gop_gcs

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/joshcarp/gop/app"
	"github.com/stretchr/testify/require"
)

func TestGCSCache(t *testing.T) {
	r := New(app.AppConfig{CacheLocation: "bucket"})
	fakeGCS := fakeServer{
		s: fakestorage.NewServer([]fakestorage.Object{}),
	}
	fakeGCS.s.CreateBucketWithOpts(fakestorage.CreateBucketOpts{Name: r.AppConfig.CacheLocation})
	r.upload = fakeGCS.uploadinMem
	req := app.NewObject("github.com/anz-bank/sysl/tests/bananatree.sysl", "e78f4afc524ad8d1a1a4740779731d706b7b079b")
	req.Content = []byte(bananatree)
	require.NoError(t, r.Cache(*req))
	obj, err := fakeGCS.s.GetObject(r.AppConfig.CacheLocation, "github.com/anz-bank/sysl/tests/bananatree.sysl@e78f4afc524ad8d1a1a4740779731d706b7b079b")
	require.NoError(t, err)
	require.Equal(t, req.Content, obj.Content)
}

type fakeServer struct {
	s *fakestorage.Server
}

func (s *fakeServer) uploadinMem(bucket string, object string, r io.Reader) error {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	s.s.CreateObject(fakestorage.Object{
		BucketName: bucket,
		Name:       object,
		Content:    bytes,
	})
	return nil
}

const bananatree = `Bananatree [package="bananatree"]:
  !type Banana:
    id <: int
    title <: string

  /banana:
    /{id<:int}:
      GET:
        return Banana

  /morebanana:
    /{id<:int}:
      GET:
        return Banana
`
