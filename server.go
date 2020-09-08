package gop

import (
	"context"

	gop2 "github.com/joshcarp/gop/gop"

	"github.com/anz-bank/pkg/log"

	"github.com/joshcarp/gop/app"

	"github.com/joshcarp/gop/gen/pkg/servers/gop"
)

/* Server is a struct that has a logger and a gopper*/
type Server struct {
	gop2.Gopper
}

/* Get implements the get endpoint*/
func (s *Server) Get(ctx context.Context, req *gop.GetRequest, client gop.GetClient) (*gop.Object, error) {
	var res gop.Object
	var cached bool
	var err error
	repo, resource, version, err := app.ProcessRequest(req.Resource)
	if err != nil {
		log.Info(ctx, err)
		return nil, err
	}
	res, cached, err = s.Retrieve(repo, resource, version)
	if err != nil || res.Content == nil || len(res.Content) == 0 {
		log.Info(ctx, "Resource not found", err)
		return nil, err
	}
	if !cached {
		if err := s.Cache(res); err != nil {
			log.Info(ctx, "Resource not found", err)
		}
	}
	return &res, nil
}
