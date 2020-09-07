package cacher_gcs

import (
	"context"
	"io"
	"testing"
	"time"

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
	req.Processed = &bananatreepbjson
	req.Content = bananatree
	require.NoError(t, r.Cache(req))
	obj, err := fakeGCS.s.GetObject(r.AppConfig.CacheLocation, "github.com/anz-bank/sysl/tests/bananatree.sysl@e78f4afc524ad8d1a1a4740779731d706b7b079b")
	require.NoError(t, err)
	require.Equal(t, []byte(req.Content), obj.Content)
	obj, err = fakeGCS.s.GetObject(r.AppConfig.CacheLocation, "github.com/anz-bank/sysl/tests/bananatree.sysl.pb.json@e78f4afc524ad8d1a1a4740779731d706b7b079b")
	require.NoError(t, err)
	require.Equal(t, []byte(*req.Processed), obj.Content)
}

type fakeServer struct {
	s *fakestorage.Server
}

func (s fakeServer) uploadinMem(bucket string, object string, r io.Reader) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()
	wr := s.s.Client().Bucket(bucket).Object(object).NewWriter(ctx)
	io.Copy(wr, r)
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

var bananatreepbjson = `{"apps":{"Bananatree":{"name":{"part":["Bananatree"]}, "attrs":{"package":{"s":"bananatree"}}, "endpoints":{"GET /banana/{id}":{"name":"GET /banana/{id}", "attrs":{"patterns":{"a":{"elt":[{"s":"rest"}]}}}, "stmt":[{"ret":{"payload":"Banana"}}], "restParams":{"method":"GET", "path":"/banana/{id}", "urlParam":[{"name":"id", "type":{"primitive":"INT", "sourceContext":{"file":"temp.sysl", "start":{"line":7, "col":5}, "end":{"line":7, "col":13}}}}]}, "sourceContext":{"file":"temp.sysl", "start":{"line":8, "col":6}, "end":{"line":11, "col":2}}}, "GET /morebanana/{id}":{"name":"GET /morebanana/{id}", "attrs":{"patterns":{"a":{"elt":[{"s":"rest"}]}}}, "stmt":[{"ret":{"payload":"Banana"}}], "restParams":{"method":"GET", "path":"/morebanana/{id}", "urlParam":[{"name":"id", "type":{"primitive":"INT", "sourceContext":{"file":"temp.sysl", "start":{"line":12, "col":5}, "end":{"line":12, "col":13}}}}]}, "sourceContext":{"file":"temp.sysl", "start":{"line":13, "col":6}, "end":{"line":15}}}}, "types":{"Banana":{"tuple":{"attrDefs":{"id":{"primitive":"INT", "sourceContext":{"file":"temp.sysl", "start":{"line":3, "col":10}, "end":{"line":3, "col":10}}}, "title":{"primitive":"STRING", "sourceContext":{"file":"temp.sysl", "start":{"line":4, "col":13}, "end":{"line":4, "col":13}}}}}, "sourceContext":{"file":"temp.sysl", "start":{"line":2, "col":2}, "end":{"line":6, "col":2}}}}, "sourceContext":{"file":"temp.sysl", "start":{"line":1, "col":1}, "end":{"line":1, "col":32}}}}}`
