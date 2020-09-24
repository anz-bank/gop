package cli

import (
	"net/http/httptest"
	"testing"

	"github.com/joshcarp/gop/gop/gop_filesystem"
	"github.com/joshcarp/gop/gop/retriever/retriever_github"

	"github.com/joshcarp/gop/gop/retrievertests"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCLI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	retriever := Default(afero.NewMemMapFs(), "", "/", "", nil)
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
	githubMock := retriever_github.NewMock()
	server := httptest.NewServer(githubMock)
	defer server.Close()
	gh := retriever_github.New(nil)
	gh.Client = server.Client()
	gh.ApiBase = server.URL
	retriever := New(
		gop_filesystem.New(fs, "."),
		gop_filesystem.New(fs, "/"),
		nil,
		gh,
		nil,
		"",
		githubMock.ResolveHash)
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
	githubMock := retriever_github.NewMock()
	server := httptest.NewServer(githubMock)
	defer server.Close()
	gh := retriever_github.New(nil)
	gh.Client = server.Client()
	gh.ApiBase = server.URL
	retriever := New(
		gop_filesystem.New(fs, "."),
		gop_filesystem.New(fs, "/"),
		nil,
		gh,
		nil,
		"test.mod",
		githubMock.ResolveHash)
	for resource, contents := range retrievertests.Tests {
		t.Run(resource, func(t *testing.T) {
			res, cached, err := retriever.Retrieve(resource)
			require.NoError(t, err)
			require.False(t, cached)
			require.Equal(t, contents, string(res))
		})

	}
}
