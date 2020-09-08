package retriever_git

import (
	"testing"

	"github.com/joshcarp/gop/app"
	"github.com/stretchr/testify/require"
)

func TestGitRetrieve(t *testing.T) {
	r := New(app.AppConfig{})
	obj, cached, err := r.Retrieve("github.com/anz-bank/sysl", "tests/bananatree.sysl", "e78f4afc524ad8d1a1a4740779731d706b7b079b")
	require.NoError(t, err)
	require.False(t, cached)
	require.Equal(t, bananatree, obj.Content)
}

const bananatree = `Bananatree [package="bananatree"]:
  !type Banana:
    id <: int
    title <: string

  /banana:
    /{id<:int}:
      GET:
        return Banana

  /morebanana:
    /{id<:int}:
      GET:
        return Banana
`
