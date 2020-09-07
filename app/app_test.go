package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepo(t *testing.T) {
	tests := map[string][]string{
		"github.com/a/b/c.d":   {"github.com/a/b", "c.d"},
		"github.com/a/b/c/d.e": {"github.com/a/b", "c/d.e"},
	}
	for in, out := range tests {
		t.Run(in, func(t *testing.T) {
			a, b := ProcessRequest(in)
			require.Equal(t, out, []string{a, b})
		})
	}
}
