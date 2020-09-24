package gop

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ghodss/yaml"
)

func TestLoadVersion(t *testing.T) {
	//	type testcase struct {
	//		name string
	//		in   string
	//		repo string
	//		out  string
	//		diff string
	//	}
	//	tests := []testcase{
	//		{name: `simple`, in: `
	//direct:
	//- pattern: github.com/abc/def
	//  resolve: github.com/abc/def@1234
	//`, repo: `github.com/abc/def`, out: `github.com/abc/def@1234`},
	//		{name: `multiple_imports`, in: `
	//direct:
	//- pattern: github.com/abc/def@1234
	//  resolve: github.com/abc/def@1234
	//- pattern: github.com/abc/xyz@567
	//  resolve: github.com/abc/xyz@567
	//`, repo: `github.com/abc/xyz@567`, out: `github.com/abc/xyz@567`},
	//		{name: `missing_import`, in: `
	//direct:
	//- pattern: github.com/abc/def@1234
	//  resolve: github.com/abc/def@1234
	//- pattern: github.com/abc/xyz@567
	//  resolve: github.com/abc/xyz@567`, repo: `github.com/def/opo`, out: `github.com/def/opo@HEAD`,
	//			diff: `
	//- pattern: github.com/def/opo
	//  resolve: github.com/def/opo@HEAD`,
	//		},
	//	}
	//
	//	for _, e := range tests {
	//		t.Run(e.name, func(t *testing.T) {
	//			c := testGopper{
	//				contents: map[string]string{"testFile": e.in},
	//				res:      map[string]string{},
	//			}
	//			resolver := func(s string) (string, error) {
	//				a := strings.Split(s, "@")
	//				if len(a) == 2 {
	//					return a[1], nil
	//				}
	//				return "HEAD", nil
	//			}
	//			ver, err := LoadVersion(c, resolver, "testFile", e.repo)
	//			require.NoError(t, err)
	//			a, b := Modules{}, Modules{}
	//			EqualYaml(e.in+"\n"+e.diff, c.contents["testFile"], &a, &b)
	//			require.Equal(t, a, b)
	//			require.Equal(t, e.out, ver)
	//		})
	//	}
}

func EqualYaml(a, b string, i, j interface{}) {
	yaml.Unmarshal([]byte(a), i)
	yaml.Unmarshal([]byte(b), j)
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

func TestReplaceSpecificImport(t *testing.T) {
	type testcase struct {
		name     string
		content  string
		oldVer   string
		oldImp   string
		newVer   string
		newImp   string
		expected string
	}
	tests := []testcase{
		{name: `simple`, oldImp: `github.com/joshcarp`, oldVer: `ver`, newImp: `asdasd`, newVer: `1234`, content: `
github.com/joshcarp/123/123.ext@ver
`, expected: `
asdasd/123/123.ext@1234
`}, {name: `without import`, oldImp: `github.com/joshcarp`, oldVer: ``, newImp: `asdasd`, newVer: `1234`, content: `
github.com/joshcarp/123/123.ext
`, expected: "\nasdasd/123/123.ext@1234\n"}}

	for _, e := range tests {
		t.Run(e.name, func(t *testing.T) {
			a := ReplaceSpecificImport(e.content, e.oldImp, e.oldVer, e.newImp, e.newVer)
			require.Equal(t, e.expected, a)
		})
	}
}
