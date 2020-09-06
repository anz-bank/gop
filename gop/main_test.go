package gop

import (
	"context"
	"testing"

	"github.com/joshcarp/pb-mod/processor/processorsysl"
	"github.com/joshcarp/pb-mod/saver/saverfs"

	"github.com/joshcarp/pb-mod/retrieve/retrieverpbjsongit"

	"github.com/joshcarp/pb-mod/config"

	"github.com/joshcarp/pb-mod/gen/pkg/servers/pbmod"

	"github.com/stretchr/testify/require"
)

func TestRepo(t *testing.T) {
	tests := map[string][]string{
		"github.com/a/b/c.d":   {"github.com/a/b", "c.d"},
		"github.com/a/b/c/d.e": {"github.com/a/b", "c/d.e"},
	}
	for in, out := range tests {
		t.Run(in, func(t *testing.T) {
			a, b := processRequest(in)
			require.Equal(t, out, []string{a, b})
		})
	}
}

func TestGetRetrieveWithDeps(t *testing.T) {
	req := &pbmod.GetResourceListRequest{
		Resource: "github.com/anz-bank/sysl/tests/bananatree.sysl",
		Version:  "e78f4afc524ad8d1a1a4740779731d706b7b079b",
	}
	client := pbmod.GetResourceListClient{}
	a := config.AppConfig{
		SaveLocation: "/dev/null",
	}
	r := retrieverpbjsongit.RetrieveFilePBJsonGit{AppConfig: a}
	s := saverfs.SaverFs{AppConfig: a}
	p := processorsysl.ProcessorSysl{}
	serve := Server{
		Retrieve: r,
		Process:  &p,
		Save:     s,
	}
	res, err := serve.GetResource(context.Background(), req, client)
	require.NoError(t, err)
	pb := "{\"apps\":{\"Bananatree\":{\"name\":{\"part\":[\"Bananatree\"]},\"attrs\":{\"package\":{\"s\":\"bananatree\"}},\"endpoints\":{\"GET /banana/{id}\":{\"name\":\"GET /banana/{id}\",\"attrs\":{\"patterns\":{\"a\":{\"elt\":[{\"s\":\"rest\"}]}}},\"stmt\":[{\"ret\":{\"payload\":\"Banana\"}}],\"restParams\":{\"method\":\"GET\",\"path\":\"/banana/{id}\",\"urlParam\":[{\"name\":\"id\",\"type\":{\"primitive\":\"INT\",\"sourceContext\":{\"file\":\"temp.sysl\",\"start\":{\"line\":7,\"col\":5},\"end\":{\"line\":7,\"col\":13}}}}]},\"sourceContext\":{\"file\":\"temp.sysl\",\"start\":{\"line\":8,\"col\":6},\"end\":{\"line\":11,\"col\":2}}},\"GET /morebanana/{id}\":{\"name\":\"GET /morebanana/{id}\",\"attrs\":{\"patterns\":{\"a\":{\"elt\":[{\"s\":\"rest\"}]}}},\"stmt\":[{\"ret\":{\"payload\":\"Banana\"}}],\"restParams\":{\"method\":\"GET\",\"path\":\"/morebanana/{id}\",\"urlParam\":[{\"name\":\"id\",\"type\":{\"primitive\":\"INT\",\"sourceContext\":{\"file\":\"temp.sysl\",\"start\":{\"line\":12,\"col\":5},\"end\":{\"line\":12,\"col\":13}}}}]},\"sourceContext\":{\"file\":\"temp.sysl\",\"start\":{\"line\":13,\"col\":6},\"end\":{\"line\":15}}}},\"types\":{\"Banana\":{\"tuple\":{\"attrDefs\":{\"id\":{\"primitive\":\"INT\",\"sourceContext\":{\"file\":\"temp.sysl\",\"start\":{\"line\":3,\"col\":10},\"end\":{\"line\":3,\"col\":10}}},\"title\":{\"primitive\":\"STRING\",\"sourceContext\":{\"file\":\"temp.sysl\",\"start\":{\"line\":4,\"col\":13},\"end\":{\"line\":4,\"col\":13}}}}},\"sourceContext\":{\"file\":\"temp.sysl\",\"start\":{\"line\":2,\"col\":2},\"end\":{\"line\":6,\"col\":2}}}},\"sourceContext\":{\"file\":\"temp.sysl\",\"start\":{\"line\":1,\"col\":1},\"end\":{\"line\":1,\"col\":32}}}}}"
	banana := &pbmod.Object{Repo: "github.com/anz-bank/sysl", Version: "e78f4afc524ad8d1a1a4740779731d706b7b079b", Resource: "tests/bananatree.sysl", Extra: &pb, Value: "Bananatree [package=\"bananatree\"]:\n  !type Banana:\n    id \u003c: int\n    title \u003c: string\n\n  /banana:\n    /{id\u003c:int}:\n      GET:\n        return Banana\n\n  /morebanana:\n    /{id\u003c:int}:\n      GET:\n        return Banana\n"}
	require.Equal(t, *banana.Extra, *res.Extra)
	banana.Extra = nil
	res.Extra = nil
	require.NoError(t, err)
	require.Equal(t, banana, res)
}
