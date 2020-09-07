package gop

import (
	"context"
	"fmt"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type Server struct {
	Retriever
	Processor
	Cacher
}

func (s Server) GetResource(ctx context.Context, req *pbmod.GetResourceListRequest, client pbmod.GetResourceListClient) (*pbmod.Object, error) {
	var object = NewObject(req.Resource, req.Version)
	object.Version = req.Version
	if err := s.Retrieve(object); err != nil {
		return nil, err
	}
	if object.Content == "" {
		return nil, fmt.Errorf("Error loading object")
	}
	if !object.Imported {
		s.Process(object)
		s.Cache(object)
	}
	if object.Processed == nil || *object.Processed == "" {
		object.Processed = nil
	}
	return object, nil
}
