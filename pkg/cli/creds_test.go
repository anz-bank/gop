package cli

import (
	"testing"

	"github.com/anz-bank/gop/pkg/gop"

	"github.com/stretchr/testify/require"
)

func TestTokensFromString(t *testing.T) {
	type testcase struct {
		name string
		in   string
		out  map[string]string
		err  error
	}
	tests := []testcase{
		{name: "empty", in: "", err: gop.UnauthorizedError},
		{in: "github.com:1234", out: map[string]string{"github.com": "1234"}},
		{in: "gitx.com:1234", out: map[string]string{"gitx.com": "1234"}},
		{in: "github.com:1234,gitx.com:12345", out: map[string]string{"github.com": "1234", "gitx.com": "12345"}}}

	for _, e := range tests {
		t.Run(e.in+e.name, func(t *testing.T) {
			tokenmap := TokensFromString(e.in)
			require.Equal(t, e.out, tokenmap)
		})
	}
}

func TestTokensFromGitCredentialsFile(t *testing.T) {
	type testcase struct {
		name string
		in   string
		out  map[string]string
		err  error
	}
	tests := []testcase{
		{name: "empty", in: ``, out: map[string]string{}},
		{name: "simple", in: "https://user:token@github.com\n", out: map[string]string{"github.com": "token"}},
		{name: "two_entries", in: "https://user:token@github.com\nhttps://user2:token2@git2.com\n", out: map[string]string{"github.com": "token", "git2.com": "token2"}},
		{name: "ignoressh", in: "https://user:token@github.com\nhttps://user2:token2@git2.com\nssh://user:token@gitub.com", out: map[string]string{"github.com": "token", "git2.com": "token2"}}}
	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			tokenmap, err := TokensFromGitCredentialsFile([]byte(e.in))
			require.Equal(t, e.err, err)
			require.Equal(t, e.out, tokenmap)
		})
	}
}
