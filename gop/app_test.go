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

func TestIsHash(t *testing.T) {
	tests := map[string]bool{
		"26b10b9ecddaa3775cf58f19aafd48b3c24951ca": true,
		"bc8fd77eb33290d5a7c0ca66fb0de9e00018b476": true,
		"fcdc2d685bc685227f9a9fbf2924cf021f10b05c": true,
		"014984980664e9f5b03c85dc887b9af129bc364a": true,
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA": true,
		"ssssssssssssssssssssssssssssssssssssssss": false,
		"abcdefhijklmnopaosdaoisdasdasdasdasdsdss": false,
		"this/is.not/a/hash//////////////////////": false,
		"this.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa": false,
		"014984980664e9f5b03c85dc887b9af129bc364":  false,
		"014984980664e9f5":                         false,
		"0149849806":                               false,
		"feature/whatever":                         false,
		"feat":                                     false,
	}
	for in, out := range tests {
		t.Run(in, func(t *testing.T) {
			require.Equal(t, out, IsHash(in))
		})
	}
}
