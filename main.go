package main

import (
	"context"
	"log"

	retrieve2 "github.com/joshcarp/pb-mod/config"

	"github.com/joshcarp/pb-mod/processor"

	"github.com/joshcarp/pb-mod/retrieve"

	"github.com/joshcarp/pb-mod/saver"

	"github.com/joshcarp/pb-mod/server"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

func main() {
	log.Fatal(pbmod.Serve(context.Background(), LoadService))
}

func LoadService(ctx context.Context, a retrieve2.AppConfig) (*pbmod.ServiceInterface, error) {
	r := retrieve.RetrieveFilePBJsonGit{AppConfig: a}
	s := saver.SaveToFile{AppConfig: a}
	p := processor.ProcessorSysl{}
	serve := server.Server{
		Retrieve: r,
		Process:  p,
		Save:     s,
	}
	return &pbmod.ServiceInterface{
		GetResourceList: serve.GetResource,
	}, nil
}
