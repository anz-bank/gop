package retriever_git

import (
	"testing"

	"github.com/joshcarp/gop/pkg/retrievertests"

	"github.com/stretchr/testify/require"
)

/* TODO: this test is failing*/
func TestGit(t *testing.T) {
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
