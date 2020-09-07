package main

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/joshcarp/pb-mod/retrieve/retrievergit"

	"github.com/joshcarp/pb-mod/retrieve/retrievergcs"

	"github.com/joshcarp/pb-mod/saver/savergcs"

	"github.com/joshcarp/pb-mod/gop"
	"github.com/joshcarp/pb-mod/processor/processorsysl"

	"github.com/joshcarp/pb-mod/config"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

func main() {
	log.Fatal(pbmod.Serve(context.Background(), LoadService))
}

func LoadService(ctx context.Context, a config.AppConfig) (*pbmod.ServiceInterface, error) {
	r := RetrieverGitGCS{retrievergcs.RetrieverGCS{AppConfig: a, Bucketname: "gop1234"}, retrievergit.RetrieverGit{AppConfig: a}}
	s := savergcs.SaverGCS{AppConfig: a, Bucketname: "gop1234"}
	p := processorsysl.ProcessorSysl{SyslimportRegex: regexp.MustCompile(processorsysl.SyslImportRegexStr)}
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
	gcs retrievergcs.RetrieverGCS
	git retrievergit.RetrieverGit
}

func (a RetrieverGitGCS) Retrieve(res *pbmod.Object) error {
	if err := a.gcs.Retrieve(res); err != nil {
		fmt.Println(err)
	}
	if res.Value == "" {
		return a.git.Retrieve(res)
	}
	return nil
}
