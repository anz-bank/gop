package cli

import (
	"testing"

	"github.com/joshcarp/gop/gop/retrievertests"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCLI(t *testing.T) {
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
