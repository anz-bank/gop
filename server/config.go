package server

import (
	"fmt"

	"github.com/joshcarp/pb-mod/retrieve"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
)

func NewKeyValue() *pbmod.Module {
	var a string
	return &pbmod.Module{
		Extra:    &a,
		Repo:     "",
		Resource: "",
		Value:    "",
		Version:  "",
	}
}

func ImportFile(initialrepo, initialImport, ver string, retriever retrieve.Retriever) (*pbmod.Module, error) {
	var file = NewKeyValue()
	file.Repo = initialrepo
	file.Resource = initialImport
	file.Version = ver
	if err := retriever.Retriever(file); err != nil {
		return nil, err
	}
	if file.Value == "" {
		return nil, fmt.Errorf("Error loading file")
	}
	return file, nil
}
