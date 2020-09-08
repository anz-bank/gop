package gop

import (
	"net/http"
	"os"

	"github.com/joshcarp/gop/gop/retriever/retriever_git"

	"github.com/anz-bank/sysl-go/common"

	"github.com/sirupsen/logrus"

	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gen/pkg/servers/gop"
	gop2 "github.com/joshcarp/gop/gen/pkg/servers/gop"
	"github.com/joshcarp/gop/gop/gop_filesystem"
	"github.com/joshcarp/gop/gop/gop_gcs"
	"github.com/spf13/afero"
)

var fs afero.Fs

func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	server, err := ServiceHandler(app.AppConfig{
		CacheLocation: os.Getenv("CacheLocation"),
		FsType:        os.Getenv("FsType"),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	server.GetHandler(w, r)
}

type GOPPER struct {
	Gopper
	retriever_git.Retriever
}

func (a GOPPER) Retrieve(repo, resource, version string) (gop.Object, bool, error) {
	if res, cached, err := a.Gopper.Retrieve(repo, resource, version); err == nil {
		return res, cached, nil
	}
	return a.Retriever.Retrieve(repo, resource, version)
}

func NewGopper(a app.AppConfig) (*GOPPER, error) {
	r := GOPPER{}
	switch a.FsType {
	case "os":
		r.Gopper = gop_filesystem.New(afero.NewOsFs(), a)
	case "mem", "memory", "":
		if fs == nil {
			fs = afero.NewMemMapFs()
		}
		r.Gopper = gop_filesystem.New(fs, a)
	case "gcs":
		gcs := gop_gcs.New(a)
		r.Gopper = &gcs
	}
	r.Retriever = retriever_git.New(a)
	return &r, nil
}

func ServiceHandler(a app.AppConfig) (*gop2.ServiceHandler, error) {
	g, err := NewGopper(a)
	if err != nil {
		return nil, err
	}
	serve := Server{
		Logger: logrus.New(),
		Gopper: g,
	}
	handler, err := gop2.NewServiceHandler(common.DefaultCallback(), &gop2.ServiceInterface{Get: serve.Get})
	if err != nil {
		return nil, err
	}
	return handler, nil

}
