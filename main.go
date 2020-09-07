package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/joshcarp/pb-mod/gop/retriever/retriever_git"

	"github.com/joshcarp/pb-mod/gop/retriever/retriever_gcs"

	"github.com/joshcarp/pb-mod/gop/cacher/cacher_gcs"

	"github.com/joshcarp/pb-mod/gop"
	"github.com/joshcarp/pb-mod/gop/processor/processor_sysl"

	"github.com/joshcarp/pb-mod/app"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

func main() {
	log.Fatal(pbmod.Serve(context.Background(), LoadService))
}

func LoadService(ctx context.Context, a app.AppConfig) (*pbmod.ServiceInterface, error) {
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
	return &pbmod.ServiceInterface{
		GetResourceList: serve.GetResource,
	}, nil
}

type RetrieverGitGCS struct {
	gcs retriever_gcs.Retriever
	git retriever_git.Retriever
}

func (a RetrieverGitGCS) Retrieve(res *pbmod.Object) error {
	if err := a.gcs.Retrieve(res); err != nil {
		fmt.Println(err)
	}
	if res.Content == "" {
		return a.git.Retrieve(res)
	}
	return nil
}
