package gop

import (
	"context"
	"fmt"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type Server struct {
	Retriever
	Processor
	Saver
}

func (s Server) GetResource(ctx context.Context, req *pbmod.GetResourceListRequest, client pbmod.GetResourceListClient) (*pbmod.Object, error) {
	var object = NewObject(req.Resource, req.Version)
	object.Version = req.Version
	if err := s.Retrieve(object); err != nil {
		return nil, err
	}
	if object.Value == "" {
		return nil, fmt.Errorf("Error loading object")
	}
	if !object.Imported {
		s.Process(object)
		s.Save(object)
	}
	if object.Extra == nil || *object.Extra == "" {
		object.Extra = nil
	}
	return object, nil
}
