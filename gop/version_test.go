package gop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadVersion(t *testing.T) {
	type testcase struct {
		name string
		in   string
		repo string
		out  string
	}
	tests := []testcase{
		{name: `simple`, in: `github.com/abc/def@1234`, repo: "github.com/abc/def", out: `1234`},
		{name: `multiple_imports`, in: "github.com/abc/def@1234\ngithub.com/abc/xyz@567\n", repo: "github.com/abc/xyz", out: `567`},
		{name: `missing_import`, in: "github.com/abc/def@1234\n", repo: "github.com/def/xyz", out: ``},
	}

	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			ver := LoadVersion([]byte(e.in), e.repo)
			require.Equal(t, e.out, ver)
		})
	}
}

func TestResolveHash(t *testing.T) {
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
			ver, err := ResolveHash(e.in)
			require.NoError(t, err)
			require.Equal(t, e.out, ver)
		})
	}
}
