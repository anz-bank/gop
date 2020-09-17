package retriever_github

import (
	"testing"

	"github.com/joshcarp/gop/gop/retrievertests"

	"github.com/stretchr/testify/require"
)

func TestGitHub(t *testing.T) {
	retriever := New("")
	for resource, contents := range retrievertests.Tests {
		t.Run(resource, func(t *testing.T) {
			res, cached, err := retriever.Retrieve(resource)
			require.NoError(t, err)
			require.False(t, cached)
			require.Equal(t, contents, string(res))
		})
	}
}
