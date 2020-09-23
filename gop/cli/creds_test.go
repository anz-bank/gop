package cli

import (
	"testing"

	"github.com/joshcarp/gop/gop"

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
			tokenmap, err := TokensFromString(e.in)
			require.Equal(t, e.err, err)
			require.Equal(t, e.out, tokenmap)
		})
	}
}
