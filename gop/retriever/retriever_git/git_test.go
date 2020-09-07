package retriever_git

import (
	"testing"

	"github.com/joshcarp/pb-mod/app"
	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"
	"github.com/stretchr/testify/require"
)

func TestGitRetrieve(t *testing.T) {
	r := New(app.AppConfig{})
	req := app.NewObject("github.com/anz-bank/sysl/tests/bananatree.sysl", "e78f4afc524ad8d1a1a4740779731d706b7b079b")
	require.NoError(t, r.Retrieve(req))
	banana := &pbmod.Object{Repo: "github.com/anz-bank/sysl", Version: "e78f4afc524ad8d1a1a4740779731d706b7b079b", Resource: "tests/bananatree.sysl", Content: "Bananatree [package=\"bananatree\"]:\n  !type Banana:\n    id \u003c: int\n    title \u003c: string\n\n  /banana:\n    /{id\u003c:int}:\n      GET:\n        return Banana\n\n  /morebanana:\n    /{id\u003c:int}:\n      GET:\n        return Banana\n"}
	req.Processed = nil
	require.Equal(t, banana, req)
}
