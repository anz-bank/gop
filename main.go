package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joshcarp/gop/gop/gop_gcs"

	"github.com/joshcarp/gop/gop/gop_filesystem"

	"github.com/spf13/afero"

	"github.com/sirupsen/logrus"

	"github.com/joshcarp/gop/gop/retriever/retriever_git"

	gop "github.com/joshcarp/gop/gop"
	"github.com/joshcarp/gop/gop/processor/processor_sysl"

	"github.com/joshcarp/gop/app"

	gop2 "github.com/joshcarp/gop/gen/pkg/servers/gop"
)

func main() {
	log.Fatal(gop2.Serve(context.Background(), LoadService))
}

func LoadService(ctx context.Context, a app.AppConfig) (*gop2.ServiceInterface, error) {
	var gopper gop.Gopper

	switch a.FsType {
	case "os":
		gopper = gop_filesystem.New(afero.NewOsFs(), a)
	case "mem", "memory":
		gopper = gop_filesystem.New(afero.NewMemMapFs(), a)
	case "gcs":
		gopper = gop_gcs.New(a)
	}

	r := Retriever{
		primary:   gopper,
		secondary: retriever_git.New(a),
	}
	p := processor_sysl.New(a)

	serve := gop.Server{
		Logger:    logrus.New(),
		Retriever: r,
		Processor: &p,
		Cacher:    gopper,
	}
	return &gop2.ServiceInterface{
		Get: serve.Get,
	}, nil
}

type Retriever struct {
	primary   gop.Retriever
	secondary gop.Retriever
}

func (a Retriever) Retrieve(res *gop2.Object) error {
	if err := a.primary.Retrieve(res); err != nil {
		fmt.Println(err)
	}
	if res.Content == "" {
		return a.secondary.Retrieve(res)
	}
	return nil
}
