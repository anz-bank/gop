package gop

import (
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

func NewObject(resource, version string) *pbmod.Object {
	var a string
	repo, resource := processRequest(resource)
	return &pbmod.Object{
		Repo:     repo,
		Resource: resource,
		Version:  version,
		Extra:    &a,
		Value:    "",
	}
}
