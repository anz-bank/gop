package server

import (
	"context"

	"github.com/joshcarp/pb-mod/processor"

	"github.com/joshcarp/pb-mod/retrieve"

	"github.com/joshcarp/pb-mod/saver"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

type Server struct {
	Retrieve retrieve.Retriever
	Process  processor.Processor
	Save     saver.Saver
}

func (s Server) GetResource(ctx context.Context, req *pbmod.GetResourceListRequest, client pbmod.GetResourceListClient) (*pbmod.Module, error) {
	repo, resource := processRequest(req.Resource)
	files, err := ImportFile(repo, resource, req.Version, s.Retrieve)
	if !files.Imported {
		s.Process.Processor(files)
		s.Save.Saver(files)
	}
	return files, err
}
