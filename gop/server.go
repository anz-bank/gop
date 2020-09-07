package gop

import (
	"context"
	"fmt"

	"github.com/joshcarp/gop/app"

	"github.com/sirupsen/logrus"

	"github.com/joshcarp/gop/gen/pkg/servers/gop"
)

type Server struct {
	*logrus.Logger
	Retriever
	Processor
	Cacher
}

func (s *Server) Get(ctx context.Context, req *gop.GetRequest, client gop.GetClient) (*gop.Object, error) {
	var object = app.NewObject(req.Resource, req.Version)
	object.Version = req.Version
	if err := s.Retrieve(object); err != nil {
		return nil, err
	}
	if object.Content == "" {
		return nil, fmt.Errorf("Error loading object")
	}
	if !object.Imported {
		if err := s.Process(object); err != nil {
			s.Logger.Println(err)
		}
		if err := s.Cache(object); err != nil {
			s.Logger.Println(err)
		}
	}
	if object.Processed == nil || *object.Processed == "" {
		object.Processed = nil
	}
	return object, nil
}
