package gop

import (
	"context"
	"fmt"

	"github.com/joshcarp/gop/app"

	"github.com/sirupsen/logrus"

	"github.com/joshcarp/gop/gen/pkg/servers/gop"
)

/* Server is a struct that has a logger and a gopper*/
type Server struct {
	*logrus.Logger
	Gopper
}

/* Get implements the get endpoint*/
func (s *Server) Get(ctx context.Context, req *gop.GetRequest, client gop.GetClient) (*gop.Object, error) {
	var res gop.Object
	var cached bool
	var err error
	repo, resource := app.ProcessRequest(req.Resource)
	if res, cached, err = s.Retrieve(repo, resource, req.Version); err != nil {
		return nil, err
	}
	if res.Content == nil {
		s.Logger.Info("Resource not found", err)
		return nil, fmt.Errorf("Error loading object")
	}
	if !cached {
		if err := s.Cache(res); err != nil {
			s.Logger.Info("Caching failed", err)
		}
	}
	return &res, nil
}
