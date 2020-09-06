package main

import (
	"context"
	"log"
	"regexp"

	"github.com/joshcarp/pb-mod/processor/processorsysl"
	"github.com/joshcarp/pb-mod/saver/saverfs"

	"github.com/joshcarp/pb-mod/gop"

	"github.com/joshcarp/pb-mod/retrieve/retrieverpbjsongit"

	"github.com/joshcarp/pb-mod/config"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

func main() {
	log.Fatal(pbmod.Serve(context.Background(), LoadService))
}

func LoadService(ctx context.Context, a config.AppConfig) (*pbmod.ServiceInterface, error) {
	r := retrieverpbjsongit.RetrieveFilePBJsonGit{AppConfig: a}
	s := saverfs.SaverFs{AppConfig: a}
	p := processorsysl.ProcessorSysl{SyslimportRegex: regexp.MustCompile(processorsysl.SyslImportRegexStr)}
	serve := gop.Server{
		Retriever: r,
		Saver:     s,
		Processor: &p,
	}
	return &pbmod.ServiceInterface{
		GetResourceList: serve.GetResource,
	}, nil
}
