package retriever_github

import (
	"net/http/httptest"
	"testing"

	"github.com/joshcarp/gop/gop/retrievertests"

	"github.com/stretchr/testify/require"
)

func TestGitHub(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	retriever := New(nil)
	for resource, contents := range retrievertests.Tests {
		t.Run(resource, func(t *testing.T) {
			res, cached, err := retriever.Retrieve(resource)
			require.NoError(t, err)
			require.False(t, cached)
			require.Equal(t, contents, string(res))
		})
	}
}

func TestGithubMock(t *testing.T) {
	mock := GithubMock{content: retrievertests.GithubRequestPaths}
	h := httptest.NewServer(&mock)
	retriever := New(nil)
	retriever.Client = h.Client()
	retriever.ApiBase = h.URL
	defer h.Close()
	for resource, contents := range retrievertests.Tests {
		t.Run(resource, func(t *testing.T) {
			res, cached, err := retriever.Retrieve(resource)
			require.NoError(t, err)
			require.False(t, cached)
			require.Equal(t, contents, string(res))
		})
	}

}

func TestResolveHash(t *testing.T) {
	mock := NewMock()
	h := httptest.NewServer(&mock)
	retriever := New(nil)
	retriever.Client = h.Client()
	retriever.ApiBase = h.URL
	type testcase struct {
		name string
		in   string
		out  string
	}
	tests := []testcase{
		{name: "tag", in: "github.com/joshcarp/gop@test", out: "dad0c54cae43ea40f3f1b5063af680ed4521eab2"},
		{name: "branch", in: "github.com/joshcarp/gop@test2", out: "dad0c54cae43ea40f3f1b5063af680ed4521eab2"},
	}
	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			ver, err := retriever.Resolve(e.in)
			require.NoError(t, err)
			require.Equal(t, e.out, ver)
		})
	}
}
