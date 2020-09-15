package gop_filesystem

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/joshcarp/gop/app"
	"github.com/stretchr/testify/require"
)

func TestFsRetrieve(t *testing.T) {
	r := New(afero.NewMemMapFs(), app.AppConfig{})
	file, err := r.fs.Create("github.com/anz-bank/sysl/tests/bananatree.sysl@e78f4afc524ad8d1a1a4740779731d706b7b079b")
	require.NoError(t, err)
	_, err = file.Write([]byte(bananatree))
	require.NoError(t, err)
	obj, _, err := r.Retrieve("github.com/anz-bank/sysl/tests/bananatree.sysl@e78f4afc524ad8d1a1a4740779731d706b7b079b")
	require.NoError(t, err)
	require.Equal(t, []byte(bananatree), obj)
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
