package cli

import (
	"log"
	"net/http/httptest"
	"testing"

	"github.com/joshcarp/gop/pkg/modules"

	"github.com/joshcarp/gop/pkg/goppers/filesystem"
	"github.com/joshcarp/gop/pkg/retrievers/github"

	"github.com/joshcarp/gop/pkg/retrievertests"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCLI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	retriever := Default(afero.NewMemMapFs(), "", "", nil)
	for resource, contents := range retrievertests.Tests {
		t.Run(resource, func(t *testing.T) {
			res, cached, err := retriever.Retrieve(resource)
			require.NoError(t, err)
			require.False(t, cached)
			require.Equal(t, contents, string(res))
		})
	}
}

func TestCLIMock(t *testing.T) {
	fs := afero.NewMemMapFs()
	githubMock := github.NewMock()
	server := httptest.NewServer(githubMock)
	defer server.Close()
	gh := github.New(nil)
	gh.Client = server.Client()
	gh.ApiBase = server.URL
	retriever := New(
		filesystem.New(fs, "."),
		filesystem.New(fs, "/"),
		nil,
		gh,
		nil,
		nil, "", log.Printf)
	for resource, contents := range retrievertests.Tests {
		t.Run(resource, func(t *testing.T) {
			res, cached, err := retriever.Retrieve(resource)
			require.NoError(t, err)
			require.False(t, cached)
			require.Equal(t, contents, string(res))
		})

	}
}

func TestCLIMockModFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	githubMock := github.NewMock()
	server := httptest.NewServer(githubMock)
	defer server.Close()
	gh := github.New(nil)
	gh.Client = server.Client()
	gh.ApiBase = server.URL
	retriever := New(
		filesystem.New(fs, "."),
		filesystem.New(fs, "/"),
		nil,
		gh,
		nil,
		nil, "test.mod", log.Printf)
	for resource, contents := range retrievertests.Tests {
		t.Run(resource, func(t *testing.T) {
			res, cached, err := retriever.Retrieve(resource)
			require.NoError(t, err)
			require.False(t, cached)
			require.Equal(t, contents, string(res))
		})

	}
}

func TestImportReplace(t *testing.T) {
	fs := afero.NewMemMapFs()
	githubMock := github.NewMock()
	server := httptest.NewServer(githubMock)
	defer server.Close()
	gh := github.New(nil)
	gh.Client = server.Client()
	gh.ApiBase = server.URL
	retriever := New(
		filesystem.New(fs, "."),
		filesystem.New(fs, "/"),
		nil,
		modules.New(gh, "test.mod"),
		nil, //modules.NewLoader(gop_filesystem.New(fs, "/"), gh.Resolve, "test.mod"),
		nil, "test.mod", log.Printf,
	)
	for resource, contents := range retrievertests.Tests {
		t.Run(resource, func(t *testing.T) {
			res, cached, err := retriever.Retrieve(resource)
			require.NoError(t, err)
			require.False(t, cached)
			require.Equal(t, contents, string(res))
		})

	}
}
