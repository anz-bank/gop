package gop

import (
	"strings"
	"testing"

	"github.com/ghodss/yaml"

	"github.com/stretchr/testify/require"
)

func TestLoadVersion(t *testing.T) {
	type testcase struct {
		name string
		in   string
		repo string
		out  string
		diff string
	}
	tests := []testcase{
		{name: `simple`, in: `
direct:
- repo: github.com/abc/def
  hash: 1234
`, repo: `github.com/abc/def`, out: `1234`},
		{name: `multiple_imports`, in: `
direct:
- repo: github.com/abc/def@1234
  hash: 1234
- repo: github.com/abc/xyz@567
  hash: 567
`, repo: `github.com/abc/xyz@567`, out: `567`},
		{name: `missing_import`, in: `
direct:
- repo: github.com/abc/def@1234
  hash: 1234
- repo: github.com/abc/xyz@567
  hash: 567`, repo: `github.com/def/opo`, out: `github.com/def/opo@HEAD`,
			diff: `
- repo: github.com/def/opo
  hash: HEAD`,
		},
	}

	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			c := testGopper{
				contents: map[string]string{"testFile": e.in},
				res:      map[string]string{},
			}
			resolver := func(s string) (string, error) {
				a := strings.Split(s, "@")
				if len(a) == 2 {
					return a[1], nil
				}
				return "HEAD", nil
			}
			ver := LoadVersion(c, resolver, "testFile", e.repo)
			a, b := Modules{}, Modules{}
			EqualYaml(e.in+"\n"+e.diff, c.contents["testFile"], &a, &b)
			require.Equal(t, a, b)
			require.Equal(t, e.out, ver)
		})
	}
}

func EqualYaml(a, b string, i, j interface{}) {
	yaml.Unmarshal([]byte(a), i)
	yaml.Unmarshal([]byte(b), j)
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

type testGopper struct {
	contents map[string]string
	res      map[string]string
}

func (r testGopper) Retrieve(resource string) (res []byte, cached bool, err error) {
	r.res[resource] = r.contents[resource]
	return []byte(r.contents[resource]), false, nil
}

func (r testGopper) Cache(resource string, content []byte) (err error) {
	r.contents[resource] = string(content)
	return nil
}
