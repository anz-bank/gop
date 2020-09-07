package main

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/joshcarp/pb-mod/gop/retriever/retriever_git"

	"github.com/joshcarp/pb-mod/gop/retriever/retriever_gcs"

	"github.com/joshcarp/pb-mod/gop/cacher/cacher_gcs"

	"github.com/joshcarp/pb-mod/gop"
	"github.com/joshcarp/pb-mod/gop/processor/processor_sysl"

	"github.com/joshcarp/pb-mod/config"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

func main() {
	log.Fatal(pbmod.Serve(context.Background(), LoadService))
}

func LoadService(ctx context.Context, a config.AppConfig) (*pbmod.ServiceInterface, error) {
	r := RetrieverGitGCS{retriever_gcs.Retriever{AppConfig: a}, retriever_git.Retriever{AppConfig: a}}
	s := cacher_gcs.Cacher{AppConfig: a}
	p := processor_sysl.Processor{ImportRegex: regexp.MustCompile(processor_sysl.SyslImportRegexStr)}
	serve := gop.Server{
		Retriever: r,
		Cacher:    s,
		Processor: &p,
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
