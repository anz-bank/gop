package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"

	"github.com/stretchr/testify/require"
)

func TestRepo(t *testing.T) {
	tests := map[string][]string{
		"github.com/a/b/c.d":   {"https://github.com/a/b.git", "c.d"},
		"github.com/a/b/c/d.e": {"https://github.com/a/b.git", "c/d.e"},
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
		Version:  "",
	}
	client := pbmod.GetResourceClient{}
	s := &server{retrieveFile: AppConfig{}.getFromGit}
	res, _ := s.GetResource(context.Background(), req, client)
	fmt.Println(string(res.Content))
}

func TestFindImport(t *testing.T) {
	var regex = `(?:#import.*)|(?:import )(?:\/\/)?(?P<import>.*)`
	tests := map[string][]string{
		`
#import notimported
import a.sysl
import b.sysl`: {"a.sysl", "b.sysl"},
	}
	for in, out := range tests {
		t.Run(in, func(t *testing.T) {
			a := findImports(regex, strings.NewReader(in))
			require.Equal(t, out, a)
		})
	}
}
