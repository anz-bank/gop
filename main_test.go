package main

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"

	"github.com/stretchr/testify/require"
)

func TestRepo(t *testing.T) {
	tests := map[string][]string{
		"github.com/a/b/c.d":   {"github.com/a/b", "c.d"},
		"github.com/a/b/c/d.e": {"github.com/a/b", "c/d.e"},
	}
	for in, out := range tests {
		t.Run(in, func(t *testing.T) {
			a, b := processRequest(in)
			require.Equal(t, out, []string{a, b})
		})
	}
}
func TestGetRetrieve(t *testing.T) {
	req := &pbmod.GetResourceRequest{
		Resource: "github.com/anz-bank/sysl/tests/bananatree.sysl",
		Version:  "e78f4afc524ad8d1a1a4740779731d706b7b079b",
	}
	client := pbmod.GetResourceClient{}
	a := &AppConfig{saveLocation: "output"}
	s := &server{retriever: []retriever{a.retrieveFie, a.retrieveGit}, saver: a.saveToFile}
	res, err := s.GetResource(context.Background(), req, client)
	require.NoError(t, err)
	banana := []pbmod.KeyValue{{Key: "github.com/anz-bank/sysl/tests/bananatree.sysl@e78f4afc524ad8d1a1a4740779731d706b7b079b", Value: "Bananatree [package=\"bananatree\"]:\n  !type Banana:\n    id <: int\n    title <: string\n\n  /banana:\n    /{id<:int}:\n      GET:\n        return Banana\n\n  /morebanana:\n    /{id<:int}:\n      GET:\n        return Banana\n"}}
	require.Equal(t, banana, res.Content)
}

func TestGetRetrieveWithDeps(t *testing.T) {
	req := &pbmod.GetResourceRequest{
		Resource: "github.com/anz-bank/sysl/tests/model_with_deps.sysl",
		Version:  "e78f4afc524ad8d1a1a4740779731d706b7b079b",
	}
	client := pbmod.GetResourceClient{}
	a := &AppConfig{saveLocation: "output"}
	s := &server{retriever: []retriever{a.retrieveFie, a.retrieveGit}, saver: a.saveToFile}
	res, err := s.GetResource(context.Background(), req, client)
	require.NoError(t, err)
	banana := []pbmod.KeyValue{{Key: "github.com/anz-bank/sysl/tests/bananatree.sysl@e78f4afc524ad8d1a1a4740779731d706b7b079b", Value: "Bananatree [package=\"bananatree\"]:\n  !type Banana:\n    id <: int\n    title <: string\n\n  /banana:\n    /{id<:int}:\n      GET:\n        return Banana\n\n  /morebanana:\n    /{id<:int}:\n      GET:\n        return Banana\n"}}
	require.Equal(t, banana, res.Content)
}

func TestGetRetrieveWithDeps2(t *testing.T) {
	req := &pbmod.GetResourceRequest{
		Resource: "github.com/anz-bank/sysl/demo/examples/Modules/model_with_deps.sysl",
		Version:  "e78f4afc524ad8d1a1a4740779731d706b7b079b",
	}
	client := pbmod.GetResourceClient{}
	a := &AppConfig{saveLocation: "output"}
	s := &server{retriever: []retriever{a.retrieveFie, a.retrieveGit}, saver: a.saveToFile}
	res, err := s.GetResource(context.Background(), req, client)
	require.NoError(t, err)
	banana := []pbmod.KeyValue{{Key: "github.com/anz-bank/sysl/tests/bananatree.sysl@e78f4afc524ad8d1a1a4740779731d706b7b079b", Value: "Bananatree [package=\"bananatree\"]:\n  !type Banana:\n    id <: int\n    title <: string\n\n  /banana:\n    /{id<:int}:\n      GET:\n        return Banana\n\n  /morebanana:\n    /{id<:int}:\n      GET:\n        return Banana\n"}}
	require.Equal(t, banana, res.Content)
}

func TestFindImport(t *testing.T) {
	tests := map[string][]string{
		`
#import notimported
import a.sysl
import b.sysl`: {"a.sysl", "b.sysl"},
	}
	for in, out := range tests {
		t.Run(in, func(t *testing.T) {
			a := findImports(syslimportRegex, []byte(in))
			require.Equal(t, out, a)
		})
	}
}

func TestDoImport(t *testing.T) {
	resources := map[string]string{
		`a.sysl`: `import b.sysl
import d.sysl

Appa:
	endpoint:
		...`,
		`b.sysl`: `
import c.sysl
Appb:
	...`,
		`c.sysl`: `
Appc:
	endpoint:
		...`, `d.sysl`: `Appd:
	endpoint:
		...`,
	}
	content, err := doImport("", `a.sysl`, "", save, retrieveFromMap, tester{resources: resources}.importerTest)
	require.NoError(t, err)
	for i, e := range content {
		require.Equal(t, resources[strings.TrimRight(strings.TrimLeft(i, "/"), "@")], e)
	}
}

type tester struct {
	resources map[string]string
}

func (t tester) importerTest(repo, resource, version string) (contents io.Reader, err error) {
	cont, ok := t.resources[resource]
	if !ok {
		return nil, fmt.Errorf("oh no")
	}
	return strings.NewReader(cont), nil
}
