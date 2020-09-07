package gop

import (
	"context"
	"fmt"

	"github.com/joshcarp/pb-mod/app"

	"github.com/sirupsen/logrus"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type Server struct {
	*logrus.Logger
	Retriever
	Processor
	Cacher
}

func (s *Server) Get(ctx context.Context, req *pbmod.GetRequest, client pbmod.GetClient) (*pbmod.Object, error) {
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
