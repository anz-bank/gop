package main

import (
	"context"
	"fmt"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type server struct {
	retrievers []retriever
	savers     []saver
	posts      []post
	AppConfig
}

func (s server) GetResource(ctx context.Context, req *pbmod.GetResourceListRequest, client pbmod.GetResourceListClient) (*pbmod.KeyValue, error) {
	repo, resource := processRequest(req.Resource)
	files, err := importFile(repo, resource, req.Version, s.retrievers)
	if !files.Imported {
		for _, p := range s.posts {
			p(files)
		}
		for _, s := range s.savers {
			if err := s(files); err != nil {
				fmt.Println(err)
			}
		}
	}
	return files, err
}
