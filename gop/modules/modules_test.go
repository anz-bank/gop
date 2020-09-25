package modules

import (
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
			a, _, err := RetrieveAndReplace(retr, i.resource, i.importFile)
			require.NoError(t, err)
			require.Equal(t, i.out, string(a))
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
