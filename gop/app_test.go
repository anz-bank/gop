package gop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessRequest(t *testing.T) {
	tests := map[string][]string{
		"github.com/a/b/c.d@123":    {"github.com/a/b", "c.d", "123"},
		"github.com/a/b/c/d.e@123":  {"github.com/a/b", "c/d.e", "123"},
		"github.com/a/b/c/d.e@1234": {"github.com/a/b", "c/d.e", "1234"},
		"github.com/a/b/c/d.e":      {"", "github.com/a/b/c/d.e", ""},
		"github.com/a/b@d":          {"github.com/a/b", "", "d"},
		"resource.ext":              {"", "resource.ext", ""},
	}
	for in, out := range tests {
		t.Run(in, func(t *testing.T) {
			a, b, c, err := ProcessRequest(in)
			require.NoError(t, err)
			require.Equal(t, out, []string{a, b, c})
		})
	}
}

func TestProcessRepo(t *testing.T) {
	tests := map[string][]string{
		"github.com/a/b/c.d@123":    {"github.com/a/b", "c.d", "123"},
		"github.com/a/b/c/d.e@123":  {"github.com/a/b", "c/d.e", "123"},
		"github.com/a/b/c/d.e@1234": {"github.com/a/b", "c/d.e", "1234"},
		"github.com/a/b/c/d.e":      {"github.com/a/b", "c/d.e", ""},
		"github.com/a/b@d":          {"github.com/a/b", "", "d"},
		"github.com":                {"github.com", "", ""},
	}
	for in, out := range tests {
		t.Run(in, func(t *testing.T) {
			a, b, c := ProcessRepo(in)
			require.Equal(t, out, []string{a, b, c})
		})
	}
}
