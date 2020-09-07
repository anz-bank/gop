package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/joshcarp/gop/gop/retriever/retriever_git"

	"github.com/joshcarp/gop/gop/retriever/retriever_gcs"

	"github.com/joshcarp/gop/gop/cacher/cacher_gcs"

	"github.com/joshcarp/gop/gop"
	"github.com/joshcarp/gop/gop/processor/processor_sysl"

	"github.com/joshcarp/gop/app"

	"github.com/joshcarp/gop/gen/pkg/servers/gop"
)

func main() {
	log.Fatal(gop.Serve(context.Background(), LoadService))
}

func LoadService(ctx context.Context, a app.AppConfig) (*gop.ServiceInterface, error) {
	r := RetrieverGitGCS{
		gcs: retriever_gcs.New(a),
		git: retriever_git.New(a),
	}
	p := processor_sysl.New(a)
	c := cacher_gcs.New(a)

	serve := gop.Server{
		Logger:    logrus.New(),
		Retriever: r,
		Processor: &p,
		Cacher:    c,
	}
	return &gop.ServiceInterface{
		Get: serve.Get,
	}, nil
}

type RetrieverGitGCS struct {
	gcs retriever_gcs.Retriever
	git retriever_git.Retriever
}

func (a RetrieverGitGCS) Retrieve(res *gop.Object) error {
	if err := a.gcs.Retrieve(res); err != nil {
		fmt.Println(err)
	}
	if res.Content == "" {
		return a.git.Retrieve(res)
	}
	return nil
}
