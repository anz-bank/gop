package modules

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/joshcarp/gop/gop/retrievertests"
)

func TestRetrieveAndReplace(t *testing.T) {
	type testcases struct {
		name       string
		resource   string
		importFile string
		out        string
		files      map[string]string
	}
	var tests = []testcases{
		{
			name:       "simple",
			resource:   "github.com/user/repo/file.ext@ver",
			importFile: "gop_modules/gop.yaml",
			out:        "import github.com/abc/def@1234",
			files: map[string]string{
				"github.com/user/repo/gop_modules/gop.yaml@ver": "imports:\n    github.com/abc/def: github.com/abc/def@1234",
				"github.com/user/repo/file.ext@ver":             "import github.com/abc/def",
			},
		},
		{
			name:       "two imports in mod file",
			resource:   "github.com/user/repo/file.ext@ver",
			importFile: "gop_modules/gop.yaml",
			out:        "import github.com/xyz/def@567",
			files: map[string]string{
				"github.com/user/repo/gop_modules/gop.yaml@ver": "imports:\n    github.com/abc/def: github.com/abc/def@1234\n    github.com/xyz/def: github.com/xyz/def@567",
				"github.com/user/repo/file.ext@ver":             "import github.com/xyz/def",
			},
		},
		{
			name:       "two imports",
			resource:   "github.com/user/repo/file.ext@ver",
			importFile: "gop_modules/gop.yaml",
			out:        "import github.com/xyz/def@567\nimport github.com/abc/def@1234",
			files: map[string]string{
				"github.com/user/repo/gop_modules/gop.yaml@ver": "imports:\n    github.com/abc/def: github.com/abc/def@1234\n    github.com/xyz/def: github.com/xyz/def@567",
				"github.com/user/repo/file.ext@ver":             "import github.com/xyz/def\nimport github.com/abc/def",
			},
		},
		{
			name:       "missing import",
			resource:   "github.com/user/repo/file.ext@ver",
			importFile: "gop_modules/gop.yaml",
			out:        "import github.com/xyz/def\nimport github.com/abc/def@1234",
			files: map[string]string{
				"github.com/user/repo/gop_modules/gop.yaml@ver": "imports:\n    github.com/abc/def: github.com/abc/def@1234",
				"github.com/user/repo/file.ext@ver":             "import github.com/xyz/def\nimport github.com/abc/def",
			},
		},
	}
	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			retr := retrievertests.New(i.files)
			a, _, err := New(retr, i.importFile).Retrieve(i.resource)
			require.NoError(t, err)
			require.Equal(t, i.out, string(a))
		})
	}
}

func TestUpdate(t *testing.T) {
	type testcases struct {
		name        string
		pattern     string
		oldResolved string
		new         string
		out         string
		importFile  string
		files       map[string]string
	}
	tests := []testcases{{
		name:        "simple",
		pattern:     "github.com/abc/def",
		oldResolved: "github.com/abc/def@1234",
		new:         "github.com/abc/def@567",
		out:         "github.com/abc/def@567",
		importFile:  "gop_modules/gop.yaml",
		files: map[string]string{
			"gop_modules/gop.yaml": "imports:\n    github.com/abc/def: github.com/abc/def@1234",
		},
	}, {
		name:        "more than one import",
		pattern:     "github.com/abc/def",
		oldResolved: "github.com/abc/def@1234",
		new:         "github.com/abc/def@567",
		out:         "github.com/abc/def@567",
		importFile:  "gop_modules/gop.yaml",
		files: map[string]string{
			"gop_modules/gop.yaml": "imports:\n    github.com/abc/def: github.com/abc/def@1234\n    github.com/xyz/xyz: github.com/xyz/xyz@xyz",
		},
	},
		{
			name:        "same repo imported under different patterns",
			pattern:     "github.com/abc/def",
			oldResolved: "github.com/abc/def@1234",
			new:         "github.com/abc/def@567",
			out:         "github.com/abc/def@567",
			importFile:  "gop_modules/gop.yaml",
			files: map[string]string{
				"gop_modules/gop.yaml": "imports:\n    github.com/abc/def@123: github.com/abc/def@1234\n    github.com/abc/def: github.com/abc/def@1234",
			},
		}, {
			name:        "resolve to HEAD",
			pattern:     "github.com/abc/def",
			oldResolved: "github.com/abc/def@1234",
			new:         "github.com/abc/def",
			out:         "github.com/abc/def@HEAD",
			importFile:  "gop_modules/gop.yaml",
			files: map[string]string{
				"gop_modules/gop.yaml": "imports:\n    github.com/abc/def@123: github.com/abc/def@1234\n    github.com/abc/def: github.com/abc/def@1234",
			},
		}}
	f := func(s string) (string, error) {
		if !strings.Contains(s, "@") {
			return s + "@HEAD", nil
		}
		return s, nil
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			retr := retrievertests.New(test.files)
			b := NewLoader(retr, f, test.importFile)
			resolved := b.Resolve(test.pattern)
			require.Equal(t, test.oldResolved, resolved)
			err := b.Update(test.pattern, test.new)
			require.NoError(t, err)
			resolved = b.Resolve(test.pattern)
			require.Equal(t, test.out, resolved)
		})
	}
}

func TestReplaceImports(t *testing.T) {
	type testcases struct {
		name    string
		modfile string
		file    string
		out     string
	}
	var tests = []testcases{
		{
			name:    "simple",
			file:    "import github.com/abc/def",
			modfile: "imports:\n    github.com/abc/def: github.com/abc/def@1234",
			out:     "import github.com/abc/def@1234",
		},
		{
			name:    "two imports in mod file",
			file:    "import github.com/xyz/def",
			modfile: "imports:\n    github.com/abc/def: github.com/abc/def@1234\n    github.com/xyz/def: github.com/xyz/def@567",
			out:     "import github.com/xyz/def@567",
		},
		{
			name:    "two imports",
			file:    "import github.com/xyz/def\nimport github.com/abc/def",
			modfile: "imports:\n    github.com/abc/def: github.com/abc/def@1234\n    github.com/xyz/def: github.com/xyz/def@567",
			out:     "import github.com/xyz/def@567\nimport github.com/abc/def@1234",
		},
		{
			name:    "missing import",
			file:    "import github.com/xyz/def\nimport github.com/abc/def",
			modfile: "imports:\n    github.com/abc/def: github.com/abc/def@1234",
			out:     "import github.com/xyz/def\nimport github.com/abc/def@1234",
		},
	}
	for _, i := range tests {
		t.Run(i.name, func(t *testing.T) {
			b, err := ReplaceImports([]byte(i.modfile), []byte(i.file))
			require.NoError(t, err)
			require.Equal(t, i.out, string(b))
		})
	}
}
