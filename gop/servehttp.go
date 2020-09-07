package gop

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joshcarp/gop/gop/retriever/retriever_git"

	"github.com/anz-bank/sysl-go/common"

	"github.com/joshcarp/gop/gop/processor/processor_sysl"

	"github.com/sirupsen/logrus"

	"github.com/joshcarp/gop/app"
	gop2 "github.com/joshcarp/gop/gen/pkg/servers/gop"
	"github.com/joshcarp/gop/gop/gop_filesystem"
	"github.com/joshcarp/gop/gop/gop_gcs"
	"github.com/spf13/afero"
)

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
	git Retriever
	Processor
}

func (a GOPPER) Retrieve(res *gop2.Object) error {
	if err := a.Gopper.Retrieve(res); err != nil {
		fmt.Println(err)
	}
	if res.Content == "" {
		return a.git.Retrieve(res)
	}
	return nil
}

func NewGopper(a app.AppConfig) (*GOPPER, error) {
	r := GOPPER{}
	switch a.FsType {
	case "os":
		r.Gopper = gop_filesystem.New(afero.NewOsFs(), a)
	case "mem", "memory", "":
		r.Gopper = gop_filesystem.New(afero.NewMemMapFs(), a)
	case "gcs":
		gcs := gop_gcs.New(a)
		r.Gopper = &gcs
	}
	r.git = retriever_git.New(a)
	r.Processor = processor_sysl.New(a)
	return &r, nil
}

func ServiceHandler(a app.AppConfig) (*gop2.ServiceHandler, error) {
	g, err := NewGopper(a)
	if err != nil {
		return nil, err
	}
	serve := Server{
		Logger:    logrus.New(),
		Retriever: g,
		Processor: g,
		Cacher:    g,
	}
	handler, err := gop2.NewServiceHandler(common.DefaultCallback(), &gop2.ServiceInterface{Get: serve.Get})
	if err != nil {
		return nil, err
	}
	return handler, nil

}
